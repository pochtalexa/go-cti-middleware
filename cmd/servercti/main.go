package main

import (
	"github.com/pochtalexa/go-cti-middleware/internal/server/api"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO тесты
	// TODO авторизация агента + DB
	// TODO sync.Mutex
	// TODO Обработка ошибок
	// TODO обработка ответа CTI на отправленные команды

	appConfig := config.NewConfig()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := appConfig.ReadConfigFile(); err != nil {
		log.Fatal().Err(err).Msg("ReadConfigFile")
	}
	if appConfig.Settings.LogLevel != "info" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	wsConn, err := cti.Init(appConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("cti.Init")
	}
	defer wsConn.Close()

	go ws.ReadMessage(wsConn, storage.AgentsInfo)

	if err := cti.InitCTISess(wsConn); err != nil {
		log.Fatal().Err(err).Msg("InitCTISess")
	}

	if err := cti.AttachUser(wsConn, "agent"); err != nil {
		log.Fatal().Err(err).Msg("AttachUser")
	}

	uHTTP := appConfig.HttpAPI.Host + ":" + appConfig.HttpAPI.Port
	if err := api.RunAPI(uHTTP); err != nil {
		log.Fatal().Err(err).Msg("RunAPI")
	}
}
