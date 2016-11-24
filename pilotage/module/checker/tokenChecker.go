package checker

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"
)

type tokenChecker struct {
}

func (checker *tokenChecker) support(eventMap map[string]string) bool {
	return true
}

func (checker *tokenChecker) Check(eventMap map[string]string, expectedToken string, reqHeader http.Header, reqBody []byte) (bool, error) {
	passCheck := false

	switch eventMap["sourceType"] {
	case "customize":
		if expectedToken == eventMap["token"] {
			passCheck = true
		}
	case "gitlab":
		if expectedToken == eventMap["token"] {
			passCheck = true
		}
	case "github":
		mac := hmac.New(sha1.New, []byte(eventMap["token"]))
		mac.Write(reqBody)
		expectedMAC := mac.Sum(nil)
		expectedSig := "sha1=" + hex.EncodeToString(expectedMAC)

		if expectedSig == eventMap["token"] {
			passCheck = true
		}
	}
	return passCheck, nil
}
