package checker

import (
	"errors"
	"net/http"

	"github.com/containerops/configure"
)

type Checker interface {
	Support(eventMap map[string]string) bool
	Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error)
}

func GetWorkflowExecCheckerList() ([]Checker, error) {
	checkers := make([]Checker, 0)
	checkerNameList := configure.GetStringSlice("auth.checker")

	for _, checkerName := range checkerNameList {
		checker, err := getChecker(checkerName)
		if err != nil {
			return nil, err
		}
		checkers = append(checkers, checker)
	}

	return checkers, nil
}

func getChecker(checkName string) (Checker, error) {
	switch checkName {
	case "tokenChecker":
		return new(tokenChecker), nil
	case "strChecker":
		return new(strChecker), nil
	}
	return nil, errors.New("error when getChecker by checkerName:" + checkName + ", checker not found")
}
