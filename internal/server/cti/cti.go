package cti

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog/log"
)

func InitCTISess(c *websocket.Conn) error {

	messInitConn := storage.NewWsCommand()
	messInitConn.Name = "SetProtocolVersion"
	messInitConn.ProtocolVersion = "13"

	if err := ws.SendCommand(c, messInitConn); err != nil {
		return fmt.Errorf("initCTISess: %w", err)
	}

	log.Info().Msg("InitCTISess - success")
	return nil
}

func AttachUser(c *websocket.Conn, login string) error {

	messAttachUser := storage.NewWsCommand()
	messAttachUser.Name = "AttachToUser"
	messAttachUser.Login = login

	if err := ws.SendCommand(c, messAttachUser); err != nil {
		return fmt.Errorf("AttachUser: %w", err)
	}

	log.Info().Str("login", login).Msg("AttachUser - success")
	return nil
}
