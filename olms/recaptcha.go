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

func challenge(action, remoteip, response string) (bool, error) {
	data := url.Values{"secret": {SecretKey}, "remoteip": {remoteip}, "response": {response}}
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
	resp, err := client.PostForm("https://www.recaptcha.net/recaptcha/api/siteverify", data)
	if err != nil {
		log.Printf("Post error: %s\n", err)
		return false, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read error: could not read body: %v", err)
		return false, err
	}
	var r reCAPTCHA
	if err := json.Unmarshal(body, &r); err != nil {
		log.Printf("Read error: got invalid JSON: %v", err)
		return false, err
	}
	if !r.Success || r.Score < 0.3 || action != r.Action {
		log.Println("Challenge failed.", r)
		return false, nil
	}
	return true, nil
}
