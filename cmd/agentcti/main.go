package main

import (
	"encoding/json"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/httpconf"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/pgui"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/storage"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func main() {
	resp := make(map[string]interface{})

	httpconf.Init()
	flags.ParseFlags()

	go pgui.Init()

	req, _ := http.NewRequest(http.MethodGet, flags.ServAddr, nil)

	for range time.Tick(time.Second * 1) {

		res, err := httpconf.HTTPClient.Do(req)
		if err != nil {
			log.Fatal().Err(err).Msg("httpClient.Do")
		}
		defer res.Body.Close()

		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(&resp); err != nil {
			log.Fatal().Err(err).Msg("Decode")
		}

		storage.AgentEvents.UserState = resp

		result, _ := storage.AgentEvents.ToString("UserState")

		pgui.UserState.SetText(result)

		//log.Info().Str("resp", fmt.Sprintln(resp)).Msg("")
		//log.Info().Str("resp[state]", fmt.Sprintln(resp["state"])).Msg("")

	}

}
