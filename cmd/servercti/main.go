package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"sync"
)

var (
	AgentInfo = storage.NewAgentInfo()
)

func wsSendCommand(c *websocket.Conn, wsMessage *storage.WsCommand) error {

	body, err := json.Marshal(wsMessage)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	err = c.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return fmt.Errorf("wsSendMessage: %w", err)
	}

	return nil
}

// горутина
func wsReadMessage(c *websocket.Conn) {
	wsEvent := storage.NewWsEvent()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("ReadMessage")
		}

		if err = json.Unmarshal(message, &wsEvent.Body); err != nil {
			log.Error().Err(err).Msg("Unmarshal")
		}

		wsEvent.Parse()
		if err = AgentInfo.SetEvent(wsEvent); err != nil {
			log.Error().Err(err).Msg("SetEvent")
		}

		log.Info().Str("message", fmt.Sprintln(wsEvent.Body)).Msg("wsReadMessage")
		log.Info().Str("message", wsEvent.Name).Msg("name")
		log.Info().Str("message", wsEvent.Login).Msg("login")
	}
}

func initCTISess(c *websocket.Conn) error {

	messInitConn := storage.NewWsCommand()
	messInitConn.Name = "SetProtocolVersion"
	messInitConn.ProtocolVersion = "13"

	if err := wsSendCommand(c, messInitConn); err != nil {
		return fmt.Errorf("initCTISess: %w", err)
	}

	return nil
}

func main() {
	var (
		wg sync.WaitGroup
	)
	config := config.NewConfig()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := config.ReadConfigFile(); err != nil {
		log.Fatal().Err(err).Msg("ReadConfigFile")
	}

	u := url.URL{
		Scheme: config.CtiAPI.Scheme,
		Host:   config.CtiAPI.Host + ":" + config.CtiAPI.Port,
		Path:   config.CtiAPI.Path}
	log.Info().Str("ws connecting to", u.String()).Msg("")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("dial")
	}
	defer c.Close()
	log.Info().Str("ws connected", u.String()).Msg("")

	wg.Add(1)
	go wsReadMessage(c)

	if err := initCTISess(c); err != nil {
		log.Fatal().Err(err).Msg("initCTISess")
	}

	messAttachUser := storage.NewWsCommand()
	messAttachUser.Name = "AttachToUser"
	messAttachUser.Login = "agent"

	if err = wsSendCommand(c, messAttachUser); err != nil {
		log.Fatal().Err(err).Msg("messAttachUser")
	}

	wg.Wait()
}
