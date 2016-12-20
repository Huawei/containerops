package checker

import (
	"errors"
	"net/http"
	"strings"
)

type strChecker struct {
}

func (checker *strChecker) Support(eventMap map[string]string) bool {
	if eventMap["sourceType"] == "gitlab" || eventMap["sourceType"] == "github" {
		return true
	}

	return false
}

func (checker *strChecker) Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error) {
	passCheck := false

	if eventMap["sourceType"] != "gitlab" && eventMap["sourceType"] != "github" {
		passCheck = false
		return passCheck, errors.New("strChecker doesn't support type:" + eventMap["sourceType"])
	}

	switch eventMap["eventName"] {
	case "Merge Request Hook":
		if strings.Contains(string(reqBody), `"action":"open"`) {
			passCheck = true
		}
	case "pull_request":
		if strings.Contains(string(reqBody), `"action":"opened"`) {
			passCheck = true
		}
	}

	if passCheck {
		return passCheck, nil
	}

	return passCheck, errors.New("strChecker failed")
}
