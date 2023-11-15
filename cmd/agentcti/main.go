package main

import (
	"encoding/json"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/flags"
	"github.com/pochtalexa/go-cti-middleware/internal/agent/pgui"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

func main() {
	httpClient := http.Client{}
	resp := make(map[string]interface{})

	flags.ParseFlags()

	go pgui.Init()

	uServer := "http://localhost:9595/"

	req, _ := http.NewRequest(http.MethodGet, uServer, nil)

	for range time.Tick(time.Second * 1) {
		userState := make(map[string]interface{})

		res, err := httpClient.Do(req)
		if err != nil {
			log.Fatal().Err(err).Msg("httpClient.Do")
		}
		defer res.Body.Close()

		dec := json.NewDecoder(res.Body)
		if err := dec.Decode(&resp); err != nil {
			log.Fatal().Err(err).Msg("Decode")
		}

		userState["state"] = resp["state"]
		userState["substates"] = resp["substates"]
		userState["time"] = resp["time"]
		userState["reason"] = resp["reason"]

		//pgui.UserState.SetText(fmt.Sprintln(userState))

		pgui.UserState.SetText("TEST")

		//log.Info().Str("resp", fmt.Sprintln(resp)).Msg("")
		//log.Info().Str("resp[state]", fmt.Sprintln(resp["state"])).Msg("")

	}

}
