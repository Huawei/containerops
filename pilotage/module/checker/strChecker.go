package checker

import (
	"errors"
	"net/http"
	"strings"
)

type strChecker struct {
}

func (checker *strChecker) Support(eventMap map[string]string) bool {
	if eventMap["sourceType"] == "gitlab" {
		return true
	}

	return false
}

func (checker *strChecker) Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error) {
	passCheck := false

	if eventMap["sourceType"] != "gitlab" {
		passCheck = false
		return passCheck, errors.New("strChecker doesn't support type:" + eventMap["sourceType"])
	}

	switch eventMap["eventName"] {
	case "Merge Request Hook":
		if strings.Contains(string(reqBody), `"action":"open"`) {
			passCheck = true
		}
	}

	return passCheck, nil
}
