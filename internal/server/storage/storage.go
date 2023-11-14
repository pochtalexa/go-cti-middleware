package storage

import (
	"fmt"
)

// AgentEvents возможные события относительно оператора
type AgentEvents struct {
	UserState              map[string]interface{}
	NewCall                map[string]interface{}
	CallStatus             map[string]interface{}
	CurrentCall            map[string]interface{}
	CallParams             map[string]interface{}
	Calls                  map[string]interface{}
	OnClose                map[string]interface{}
	OnTransferCall         map[string]interface{}
	TransferSucceed        map[string]interface{}
	TransferFailed         map[string]interface{}
	TransferCallReturned   map[string]interface{}
	SetSessionModeResponse map[string]interface{}
	CallParamsUpdated      map[string]interface{}
	LocalParamsUpdated     map[string]interface{}
	Conferences            map[string]interface{}
	Ok                     map[string]interface{}
	Error                  map[string]interface{}
	ParseError             map[string]interface{}
}

// AgentsInfo мапа с ключом - логин оператора
type AgentsInfo struct {
	Events map[string]AgentEvents
}

// WsCommand команда в сторону CTI API
type WsCommand struct {
	Name            string `json:"name"` // название команды или события
	Login           string `json:"login,omitempty"`
	Rid             string `json:"rid,omitempty"` // Значение данного поля ответного сообщения соответствует значению этого же поля команды. Если событие не является ответом на команду, то поле отсутствует.
	Cid             string `json:"cid,omitempty"` // Идентификатор обращения (call identifier) в контексте оператора. Указывается как идентификатор обращения во всех сообщениях, которые касаются обработки обращения.
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	ParamName       string `json:"paramName,omitempty"`
	ParamValue      string `json:"paramValue,omitempty"`
	On              bool   `json:"on,omitempty"`
	Enable          bool   `json:"enable,omitempty"`
	Target          string `json:"target,omitempty"`
	DTMFString      string `json:"DTMFString,omitempty"`
	Url             string `json:"url,omitempty"`
}

// WsEvent событие или ответ от CTI API
type WsEvent struct {
	Name       string
	Login      string
	Body       map[string]interface{}
	ErrorNames []string
}

func NewWsCommand() *WsCommand {
	return &WsCommand{}
}

func NewWsEvent() *WsEvent {
	return &WsEvent{
		ErrorNames: []string{"Error", "ParseError"},
	}
}

func NewAgentsInfo() *AgentsInfo {
	return &AgentsInfo{
		Events: make(map[string]AgentEvents),
	}
}

func (w *WsEvent) Parse() {
	if _, ok := w.Body["name"]; ok {
		w.Name = fmt.Sprintf("%v", w.Body["name"])
	} else {
		w.Name = ""
	}

	if _, ok := w.Body["login"]; ok {
		w.Login = fmt.Sprintf("%v", w.Body["login"])
	} else {
		w.Login = ""
	}
}

func (a *AgentsInfo) SetEvent(event *WsEvent) error {
	// сохраняем текущие события по оператору и обновляем
	curEvents := a.Events[event.Login]

	switch event.Name {
	case "UserState":
		curEvents.UserState = event.Body
	case "NewCall":
		curEvents.NewCall = event.Body
	case "LocalParamsUpdated":
		curEvents.LocalParamsUpdated = event.Body
	case "CallParamsUpdated":
		curEvents.CallParamsUpdated = event.Body
	case "SetSessionModeResponse":
		curEvents.SetSessionModeResponse = event.Body
	case "TransferCallReturned":
		curEvents.TransferCallReturned = event.Body
	case "TransferFailed":
		curEvents.TransferFailed = event.Body
	case "TransferSucceed":
		curEvents.TransferSucceed = event.Body
	case "OnTransferCall":
		curEvents.OnTransferCall = event.Body
	case "OnClose":
		curEvents.OnClose = event.Body
	case "Calls":
		curEvents.Calls = event.Body
	case "CallParams":
		curEvents.CallParams = event.Body
	case "CurrentCall":
		curEvents.CurrentCall = event.Body
	case "CallStatus":
		curEvents.CallStatus = event.Body
	case "Conferences":
		curEvents.Conferences = event.Body
	case "Ok":
		curEvents.Ok = event.Body
	case "Error":
		curEvents.Error = event.Body
	case "ParseError":
		curEvents.ParseError = event.Body
	default:
		return fmt.Errorf("can not find case for key %v", event.Name)
	}

	if event.Login == "" {
		a.Events["noName"] = curEvents
	} else {
		a.Events[event.Login] = curEvents
	}

	return nil
}
