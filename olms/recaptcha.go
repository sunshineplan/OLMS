package olms

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type reCAPTCHA struct {
	Success     bool
	Score       float64
	Action      string
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string
	ErrorCodes  []string `json:"error-codes"`
}

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
	var r reCAPTCHA
	if err := json.Unmarshal(b, &r); err != nil {
		log.Println("Failed to unmarshal response:", err)
		return false
	}
	if !r.Success || r.Score < 0.5 || action != r.Action {
		log.Println("Challenge failed.", r)
		return false
	}
	return true
}
