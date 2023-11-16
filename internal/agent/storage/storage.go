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
			event := a.UserState
			result = fmt.Sprintf("state: %v, substates: %v, reason: %v",
				event["state"], event["substates"], event["reason"])
		} else {
			result = "-"
		}
	case "NewCall":
		_, ok := a.NewCall["name"]
		if ok {
			event := a.NewCall
			result = fmt.Sprintf("state: %v, direction: %v, displayName: %v",
				event["state"], event["direction"], event["displayName"])
		} else {
			result = "-"
		}
	case "CallStatus":
		_, ok := a.CallStatus["name"]
		if ok {
			event := a.CallStatus
			result = fmt.Sprintf(
				"state: %v, "+
					"params: %v, "+
					"creationTime: %v, "+
					"answerTime: %v, "+
					"hangupTime: %v, ",
				event["state"],
				event["params"],
				event["creationTime"],
				event["answerTime"],
				event["hangupTime"])
		} else {
			result = "-"
		}
	default:
		return "", fmt.Errorf("ToString - can not find event name: %s", name)
	}

	return result, nil
}
