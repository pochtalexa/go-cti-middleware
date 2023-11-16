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
	httpconf.Init()
	flags.ParseFlags()

	go pgui.Init()

	url := flags.ServAddr + "/api/v1/events/agent"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// по таймеру запрашиваем новые метрики
	for range time.Tick(time.Second * 1) {

		res, err := httpconf.HTTPClient.Do(req)
		if err != nil {
			log.Fatal().Err(err).Msg("httpClient.Do")
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNoContent {
			continue
		}

		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(&storage.AgentEvents); err != nil {
			log.Fatal().Err(err).Msg("Decode")
		}

		result, _ := storage.AgentEvents.ToString("UserState")
		pgui.UserState.SetText(result)

		result, _ = storage.AgentEvents.ToString("NewCall")
		pgui.NewCall.SetText(result)

		result, _ = storage.AgentEvents.ToString("CallStatus")
		pgui.CallStatus.SetText(result)

		//log.Info().Str("resp", fmt.Sprintln(resp)).Msg("")
		//log.Info().Str("resp[state]", fmt.Sprintln(resp["state"])).Msg("")

	}

}
