package sraw

import (
	"net/http"
	"strings"
	"time"
)

// sendAuthRequest sends an authenticated request to a url
func (a Auth) sendAuthRequest(method, url, body string) (*http.Response, error) {
	// prepare url
	if strings.ContainsAny(url, "?") {
		url += "&access_token=" + a.accToken
	} else {
		url += "?access_token=" + a.accToken
	}

	reqBody := strings.NewReader(body)
	if body == "" {
		reqBody = nil
	}

	// creates request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// sends request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
