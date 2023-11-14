package main

import (
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"sync"
)

func main() {
	var (
		wg sync.WaitGroup
	)
	appConfig := config.NewConfig()
	agentsInfo := storage.NewAgentsInfo()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := appConfig.ReadConfigFile(); err != nil {
		log.Fatal().Err(err).Msg("ReadConfigFile")
	}

	if appConfig.Settings.LogLevel != "info" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	u := url.URL{
		Scheme: appConfig.CtiAPI.Scheme,
		Host:   appConfig.CtiAPI.Host + ":" + appConfig.CtiAPI.Port,
		Path:   appConfig.CtiAPI.Path}
	log.Info().Str("ws connecting to", u.String()).Msg("")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("dial")
	}
	defer c.Close()
	log.Info().Str("ws connected", u.String()).Msg("")

	wg.Add(1)
	go ws.ReadMessage(c, agentsInfo)

	if err := cti.InitCTISess(c); err != nil {
		log.Fatal().Err(err).Msg("initCTISess")
	}

	if err := cti.AttachUser(c, "agent"); err != nil {
		log.Fatal().Err(err).Msg("AttachUser")
	}

	wg.Wait()
}
