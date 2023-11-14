package handlers

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	//reqBody := make(map[string]interface{})
	resBody := make(map[string]string)

	//dec := json.NewDecoder(r.Body)
	//if err := dec.Decode(&reqBody); err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//	log.Error().Err(err).Msg("Decode")
	//	return
	//}
	//log.Info().Str("reqBody", fmt.Sprint(reqBody)).Msg("reqBody")

	resBody["status"] = "ok"
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
