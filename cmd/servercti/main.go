package main

import (
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/api"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
)

func main() {
	// TODO тесты
	// TODO авторизация агента + DB
	// TODO sync.Mutex

	appConfig := config.NewConfig()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := appConfig.ReadConfigFile(); err != nil {
		log.Fatal().Err(err).Msg("ReadConfigFile")
	}

	if appConfig.Settings.LogLevel != "info" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	uCTI := url.URL{
		Scheme: appConfig.CtiAPI.Scheme,
		Host:   appConfig.CtiAPI.Host + ":" + appConfig.CtiAPI.Port,
		Path:   appConfig.CtiAPI.Path,
	}
	log.Info().Str("ws connecting to", uCTI.String()).Msg("")

	c, _, err := websocket.DefaultDialer.Dial(uCTI.String(), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("dial")
	}
	defer c.Close()
	log.Info().Str("ws connected", uCTI.String()).Msg("")

	go ws.ReadMessage(c, storage.AgentsInfo)

	if err := cti.InitCTISess(c); err != nil {
		log.Fatal().Err(err).Msg("initCTISess")
	}

	if err := cti.AttachUser(c, "agent"); err != nil {
		log.Fatal().Err(err).Msg("AttachUser")
	}

	uHTTP := appConfig.HttpAPI.Host + ":" + appConfig.HttpAPI.Port
	if err := api.RunAPI(uHTTP); err != nil {
		log.Fatal().Err(err).Msg("RunAPI")
	}
}
