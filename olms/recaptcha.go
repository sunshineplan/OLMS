package olms

import (
	"log"
	"net/url"
	"time"

	"github.com/sunshineplan/utils/requests"
)

type reCAPTCHA struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
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
	var r reCAPTCHA
	if err := requests.Post("https://www.recaptcha.net/recaptcha/api/siteverify", nil, data).JSON(&r); err != nil {
		log.Printf("Failed to verify response: %v", err)
		return false
	}
	if !r.Success || r.Score < 0.5 || action != r.Action {
		log.Println("Challenge failed.", r)
		return false
	}
	return true
}
