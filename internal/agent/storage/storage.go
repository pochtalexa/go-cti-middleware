package storage

import (
	"fmt"
)

var (
	AgentEvents = NewAgentEvents()
)

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

func NewAgentEvents() *StAgentEvents {
	return &StAgentEvents{}
}

//func ToString(el map[string]interface{}) string {
//	var storeList []string
//
//	for k, v := range el {
//		storeList = append(storeList, k+":"+fmt.Sprintf("%v", v))
//	}
//
//	result := strings.Join(storeList, ",")
//
//	return result
//}

func (a *StAgentEvents) ToString(name string) (string, error) {
	var result string

	switch name {
	case "UserState":
		_, ok := a.UserState["name"]
		if ok {
			result = fmt.Sprintf("state: %v, substates: %v, reason: %v",
				a.UserState["state"], a.UserState["substates"], a.UserState["reason"])
		} else {
			result = "-"
		}
	default:
		return "", fmt.Errorf("ToString - can not find event name: %s", name)
	}

	return result, nil
}
