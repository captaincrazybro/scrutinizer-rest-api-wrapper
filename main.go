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
	res, err := a.sendAuthRequest(http.MethodPost, reqURL, "")
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

func (a Auth) GetReportDetails(provType prvTp, owner string, name string) (ReportDetails, error) {
	err := a.Validate()
	if err != nil {
		return details, err
	}

	// send request
	reqURL := fmt.Sprintf("%s%s/repositories/%s/%s", Endpoint, provType.name, owner, name)
	res, err := a.sendAuthRequest(http.MethodGet, reqURL, "")
	if err != nil {
		return err
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
	var details ReportDetails{}
	err = json.Unmarshal(btsBody, &details)
	if err != nil {
		return nil, err
	}

	return &details, nil

}

type ReportDetails struct {
	Date                string  `json:"date"`
	CreatedAt           string  `json:"created_at"`
	StartDate           string  `json:"start_date"`
	EndDate             string  `json:"end_date"`
	BranchReference     string  `json:"branch_reference"`
	BaseSourceReference string  `json:"base_source_reference"`
	HeadSourceReference string  `json:"head_source_reference"`
	QualityScore        float64 `json:"quality_score"`
	QualityScoreChange  float64 `json:"quality_score_change"`
	QualityDistribution struct {
		Weights struct {
			VeryGood     float64 `json:"very_good"`
			Good         float64 `json:"good"`
			Satisfactory float64 `json:"satisfactory"`
			Pass         float64 `json:"pass"`
			Critical     float64 `json:"critical"`
		} `json:"weights"`
	} `json:"quality_distribution"`
	NbAlerts           int `json:"nb_alerts"`
	NbAlertsChange     int `json:"nb_alerts_change"`
	NbIssues           int `json:"nb_issues"`
	NbIssuesChange     int `json:"nb_issues_change"`
	TestCoverageChange int `json:"test_coverage_change"`
	NbCommits          int `json:"nb_commits"`
	NbAdditions        int `json:"nb_additions"`
	NbDeletions        int `json:"nb_deletions"`
	LargestCommits     []struct {
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Title string `json:"title"`
		Ref   string `json:"ref"`
	} `json:"largest_commits"`
	TopContributors []struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Nbcommits   int    `json:"nbCommits"`
		Nbadditions int    `json:"nbAdditions"`
		Nbdeletions int    `json:"nbDeletions"`
	} `json:"top_contributors"`
	AlgorithmChanged bool `json:"algorithm_changed"`
	Links            struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Repository struct {
			Href string `json:"href"`
		} `json:"repository"`
	} `json:"_links"`
	Embedded struct {
		Repository struct {
			Type                      string `json:"type"`
			CreatedAt                 string `json:"created_at"`
			Private                   bool   `json:"private"`
			DefaultBranch             string `json:"default_branch"`
			DevelopmentReportSettings struct {
				Enabled  bool   `json:"enabled"`
				Weekday  int    `json:"weekday"`
				Hour     int    `json:"hour"`
				Timezone string `json:"timezone"`
			} `json:"development_report_settings"`
			BranchSettings struct {
				TrackedBranches []string `json:"tracked_branches"`
			} `json:"branch_settings"`
			Login string `json:"login"`
			Name  string `json:"name"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
		} `json:"repository"`
	} `json:"_embedded"`
}