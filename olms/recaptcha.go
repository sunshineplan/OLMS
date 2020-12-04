package olms

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// SiteKey reCAPTCHA
var SiteKey string

// SecretKey reCAPTCHA
var SecretKey string

func challenge(action, remoteip string, response interface{}) bool {
	token, ok := response.(string)
	if !ok {
		return false
	}
	data := url.Values{"secret": {SecretKey}, "remoteip": {remoteip}, "response": {token}}

	resp, err := http.PostForm("https://www.recaptcha.net/recaptcha/api/siteverify", data)
	if err != nil {
		log.Println("Failed to verify response:", err)
		return false
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response:", err)
		return false
	}

	var result struct {
		Success     bool
		Score       float64
		Action      string
		ChallengeTS time.Time `json:"challenge_ts"`
		Hostname    string
		ErrorCodes  []string `json:"error-codes"`
	}
	if err := json.Unmarshal(b, &result); err != nil {
		log.Println("Failed to unmarshal response:", err)
		return false
	}
	if !result.Success || result.Score < 0.5 || action != result.Action {
		log.Println("Challenge failed.", result)
		return false
	}
	return true
}
