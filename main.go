package sraw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Auth stores authentication for requests
type Auth struct {
	// accToken - [REQUIRED] stores the access token required to create requests to scrutinizer's api
	accToken string
}

// Repository interface to store repositories
type Repository struct {
}

// prvTp stores the name of a provider type
type prvTp struct {
	name string
}

// prvdrs struct for storing providers
type prvdrs struct {
	GITHUB    prvTp
	BITBUCKET prvTp
}

// Providers stores providers
var Providers prvdrs = prvdrs{
	GITHUB:    prvTp{"g"},
	BITBUCKET: prvTp{"b"},
}

// RepositoryPayload structure for storing repository json payload
type RepositoryPayload struct {
	Type          string `json:"type"`
	CreatedAt     string `json:"created_at"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
	Login         string `json:"your-login"`
	Name          string `json:"name"`
}

// errorPayload structure for storing error json payload
type errorPayload struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

// addRepoBody structure for storing the body to pass along with the add repo request
type addRepoBody struct {
	Name         string `json:"name"`
	Org          string `json:"organization"`
	Config       string `json:"config"`
	GlobalConfig string `json:"global_config"`
}

// Validate makes sure all required parameters are set and valid
func (a Auth) Validate() error {
	if a.accToken == "" {
		return fmt.Errorf("access token needs to be set")
	}
	return nil
}

// GetRepo gets information about a scrutinizer repository and returns nil for both repo payload and error if it doesn't exist
func (a Auth) GetRepo(provType prvTp, owner string, name string) (*RepositoryPayload, error) {
	err := a.Validate()
	if err != nil {
		return nil, err
	}

	// send request
	reqURL := fmt.Sprintf("%s%s/repositories/%s/%s", Endpoint, provType.name, owner, name)
	res, err := a.sendAuthRequest(http.MethodGet, reqURL, "")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	btsBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// checks if repo exists
	errPayload := errorPayload{}
	err = json.Unmarshal(btsBody, &errPayload)
	if err != nil {
		return nil, err
	}

	if errPayload.Message != "" && errPayload.Description != "" {
		if errPayload.Message == "Not Found" {
			return nil, nil
		}
		return nil, fmt.Errorf("%s, %s", errPayload.Message, errPayload.Description)
	}

	// handles successful payload
	repoPayload := RepositoryPayload{}
	err = json.Unmarshal(btsBody, &repoPayload)
	if err != nil {
		return nil, err
	}

	return &repoPayload, nil

}

// AddRepo adds a repo to scrutinizer
func (a Auth) AddRepo(provType prvTp, owner string, name string, config string, globalConfig string) error {
	err := a.Validate()
	if err != nil {
		return err
	}

	// create body
	repoBody := addRepoBody{
		Name:         name,
		Org:          owner,
		Config:       config,
		GlobalConfig: globalConfig,
	}

	body, err := json.Marshal(repoBody)
	if err != nil {
		return err
	}

	// send request
	reqURL := fmt.Sprintf("%s%s", Endpoint, provType.name)
	res, err := a.sendAuthRequest(http.MethodPost, reqURL, string(body))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	btsBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// checks if repo exists
	errPayload := errorPayload{}
	err = json.Unmarshal(btsBody, &errPayload)
	if err != nil {
		return err
	}

	if errPayload.Message != "" && errPayload.Description != "" {
		return fmt.Errorf("%s, %s", errPayload.Message, errPayload.Description)
	}

	return nil

}
