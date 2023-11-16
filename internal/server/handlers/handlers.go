package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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

// ControlHandler принимем команду по http API и вызваем соотвествующий медот CTI API
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

// GetEventsHandler запрос на отправку текущих events
func GetEventsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO добавить автоизацию

	// проверяем что логин в запросе не пустой
	login := chi.URLParam(r, "login")
	if login == "" {
		errorText := fmt.Errorf("no login in request")
		http.Error(w, errorText.Error(), http.StatusBadRequest)
		log.Error().Err(errorText).Msg("parse url")
		return
	}

	// проверяе что есть обновленные данные
	updated, ok := storage.AgentsInfo.Updated[login]
	if !ok {
		errorText := fmt.Errorf("no key for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNotFound)
		log.Error().Err(errorText).Msg("")
		return
	}
	if !updated {
		errorText := fmt.Errorf("no updated data for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNoContent)
		log.Info().Err(errorText).Msg("")
		return
	}

	resBody, ok := storage.AgentsInfo.Events[login]
	if !ok {
		errorText := fmt.Errorf("no key for agent with login: %s", login)
		http.Error(w, errorText.Error(), http.StatusNotFound)
		log.Error().Err(errorText).Msg("")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(resBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error().Err(err).Msg("Encode")
		return
	}

	storage.AgentsInfo.Updated[login] = false

	// очищаем хранилище после отправки
	storage.AgentsInfo.DropAgentEvents(login)

	w.WriteHeader(http.StatusOK)
	log.Info().Msg("GetEventsHandler - ok")

	return
}
