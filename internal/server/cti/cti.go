package cti

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pochtalexa/go-cti-middleware/internal/server/config"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/pochtalexa/go-cti-middleware/internal/server/ws"
	"github.com/rs/zerolog/log"
	"net/url"
)

var Conn *websocket.Conn

func Init(appConfig *config.Config) (*websocket.Conn, error) {
	var err error

	uCTI := url.URL{
		Scheme: appConfig.CtiAPI.Scheme,
		Host:   appConfig.CtiAPI.Host + ":" + appConfig.CtiAPI.Port,
		Path:   appConfig.CtiAPI.Path,
	}
	log.Info().Str("ws connecting to", uCTI.String()).Msg("")

	Conn, _, err = websocket.DefaultDialer.Dial(uCTI.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial: %w", err)
	}
	log.Info().Str("ws connected", uCTI.String()).Msg("")

	return Conn, nil
}

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

func ChageStatus(c *websocket.Conn, login string, status string) error {

	messChageStatus := storage.NewWsCommand()
	messChageStatus.Name = "ChangeUserState"
	messChageStatus.Login = login
	messChageStatus.State = status

	if err := ws.SendCommand(c, messChageStatus); err != nil {
		return fmt.Errorf("ChangeUserState: %w", err)
	}

	log.Info().Str("login", login).Msg("ChangeUserState - success")
	return nil

}
