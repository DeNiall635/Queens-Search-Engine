package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/ad"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/api"
)

// Retriever is a client interface for searching the database for relevant ads
type Retriever interface {
	Get(keywords string) (*ad.Advert, error)
}

// APIAdRetrieve searches a database for relevant ads
type APIAdRetrieve struct {
	Endpoint string
	Client   http.Client
}

// Get takes keywords and returns relevant ads using the ads microservice api
func (s *APIAdRetrieve) Get(keywords string) (*ad.Advert, error) {
	resp, err := s.Client.Get(fmt.Sprintf("%s/ad?keywords=%s", s.Endpoint, url.QueryEscape(keywords)))
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
	var result ad.Advert
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
