package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pochtalexa/go-cti-middleware/internal/server/cti"
	"github.com/pochtalexa/go-cti-middleware/internal/server/storage"
	"github.com/rs/zerolog/log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	//reqBody := make(map[string]interface{})
	//resBody := make(map[string]string)
	//resBody := make(map[string]interface{})

	//dec := json.NewDecoder(r.Body)
	//if err := dec.Decode(&reqBody); err != nil {
	//	w.WriteHeader(httpconf.StatusInternalServerError)
	//	log.Error().Err(err).Msg("Decode")
	//	return
	//}
	//log.Info().Str("reqBody", fmt.Sprint(reqBody)).Msg("reqBody")

	//resBody["status"] = "ok"

	resBody := storage.AgentsInfo.Events["agent"].UserState

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error().Err(err).Msg("Encode")
		return
	}

	return
}

func ControlHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := make(map[string]string)

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error().Err(err).Msg("Decode")
		return
	}

	if err := cti.ChageStatus(cti.Conn, "agent", reqBody["ChangeUserState"]); err != nil {
		log.Error().Err(err).Msg("call ChageStatus")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Info().Str("reqBody", fmt.Sprint(reqBody)).Msg("reqBody")
}
