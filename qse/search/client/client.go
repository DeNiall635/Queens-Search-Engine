package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gitlab2.eeecs.qub.ac.uk/40173800/qse/search/api"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/search/search"
)

// Searcher is a client interface for searching the database for a query with pagination
type Searcher interface {
	Search(query string, page int) (*search.Result, error)
}

// APISearch is used to search using the search microservice API
type APISearch struct {
	Endpoint string
	Client   http.Client
}

// Search takes a query and a page and returns results, by using the search microservice API
func (s *APISearch) Search(query string, page int) (*search.Result, error) {
	resp, err := s.Client.Get(fmt.Sprintf("%s/search?q=%s&p=%d", s.Endpoint, url.QueryEscape(query), page))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var httpErr api.HTTPError
		err = json.Unmarshal(body, &httpErr)
		if err != nil {
			return nil, fmt.Errorf("Got error %d, unable to unmarshal error %v, %s", resp.StatusCode, err, body)
		}
		return nil, fmt.Errorf("Got error %d: %s", resp.StatusCode, httpErr.Message)
	}
	var result search.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
