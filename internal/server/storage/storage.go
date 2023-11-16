package storage

import (
	"fmt"
)

var AgentsInfo = NewAgentsInfo()

// TODO на перспективу использовать Redis
// TODO описать поля мап

// StAgentEvents возможные события относительно оператора
type StAgentEvents struct {
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

// StAgentsInfo мапа с ключом - логин оператора
type StAgentsInfo struct {
	Events  map[string]StAgentEvents
	Updated map[string]bool
}

// StWsCommand команда в сторону CTI API
type StWsCommand struct {
	Name            string `json:"name"` // название команды или события
	Login           string `json:"login,omitempty"`
	Rid             string `json:"rid,omitempty"` // Значение данного поля ответного сообщения соответствует значению этого же поля команды. Если событие не является ответом на команду, то поле отсутствует.
	Cid             string `json:"cid,omitempty"` // Идентификатор обращения (call identifier) в контексте оператора. Указывается как идентификатор обращения во всех сообщениях, которые касаются обработки обращения.
	ProtocolVersion string `json:"protocolVersion,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	ParamName       string `json:"paramName,omitempty"`
	ParamValue      string `json:"paramValue,omitempty"`
	State           string `json:"state,omitempty"` // Состояние программного телефона, которое необходимо установить
	On              bool   `json:"on,omitempty"`
	Enable          bool   `json:"enable,omitempty"`
	Target          string `json:"target,omitempty"`
	DTMFString      string `json:"DTMFString,omitempty"`
	Url             string `json:"url,omitempty"`
}

// StWsEvent событие или ответ от CTI API
type StWsEvent struct {
	Name       string
	Login      string
	Body       map[string]interface{}
	ErrorNames []string
}

func NewWsCommand() *StWsCommand {
	return &StWsCommand{}
}

func NewWsEvent() *StWsEvent {
	return &StWsEvent{
		ErrorNames: []string{"Error", "ParseError"},
	}
}

func NewAgentsInfo() *StAgentsInfo {
	return &StAgentsInfo{
		Events:  make(map[string]StAgentEvents),
		Updated: make(map[string]bool),
	}
}

func (w *StWsEvent) Parse() {
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

func (a *StAgentsInfo) SetEvent(event *StWsEvent) error {
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
		a.Updated["noName"] = true
	} else {
		a.Events[event.Login] = curEvents
		a.Updated[event.Login] = true
	}

	return nil
}

func (a *StAgentsInfo) DropAgentEvents(login string) {
	a.Events[login] = StAgentEvents{}
}
