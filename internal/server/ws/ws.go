package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"slices"
)

func SendCommand(c *websocket.Conn, wsMessage *storage.WsCommand) error {

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

// ReadMessage WS goroutine
func ReadMessage(c *websocket.Conn, agentsInfo *storage.AgentsInfo) {

	for {
		wsEvent := storage.NewWsEvent()
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error().Err(err).Msg("ReadMessage")
		}

		if err = json.Unmarshal(message, &wsEvent.Body); err != nil {
			log.Error().Err(err).Msg("Unmarshal")
		}

		wsEvent.Parse()

		if err = agentsInfo.SetEvent(wsEvent); err != nil {
			log.Error().Err(err).Msg("SetEvent")
		}

		if slices.Contains(wsEvent.ErrorNames, wsEvent.Name) {
			log.Error().Str("message", fmt.Sprintln(wsEvent.Body)).Msg("wsReadMessage")
		} else {
			log.Info().Str("event name", wsEvent.Name).Str("message", fmt.Sprintln(wsEvent.Body)).Msg("wsReadMessage")
		}
		//log.Info().Str("message", wsEvent.Name).Msg("name")
		//log.Info().Str("message", wsEvent.Login).Msg("login")
	}
}
