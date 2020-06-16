package ad

import (
	"encoding/json"
	"log"
	"net/http"

	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/ad"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/api"
)

const (
	keywordsParam = "keywords"
	titleParam    = "title"
	uriParam      = "uri"
	contentParam  = "content"
	keywordParam  = "keyword"
)

// Ad handles the GET /ad endpoint, providing simple parameter validation, allows retrieving HTML for an advert based on a keyword
func Ad(server ad.Server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		keywords := r.URL.Query().Get(keywordsParam)
		// Check params provided
		if keywords == "" {
			api.Error(w, "The keywords field 'keywords' is required.", http.StatusBadRequest)
			return
		}

		// Get search results
		results, err := server.Get(keywords)
		if err != nil {
			api.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(results)
		if err != nil {
			// Should not occur, panic
			log.Panic(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(json)

	})
}

// PostAd handles the POST /ad endpoint, providing simple parameter validation, allows creating an Advert
func PostAd(server ad.Server) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check params provided
		title := r.URL.Query().Get(titleParam)
		if title == "" {
			api.Error(w, "The title field 'title' is required.", http.StatusBadRequest)
			return
		}

		uri := r.URL.Query().Get(uriParam)
		if uri == "" {
			api.Error(w, "The uri field 'uri' is required.", http.StatusBadRequest)
			return
		}

		content := r.URL.Query().Get(contentParam)
		if content == "" {
			api.Error(w, "The content field 'content' is required.", http.StatusBadRequest)
			return
		}

		keyword := r.URL.Query().Get(keywordParam)
		if keyword == "" {
			api.Error(w, "The keyword field 'keyword' is required.", http.StatusBadRequest)
			return
		}

		// Get search results
		err := server.Create(title, uri, content, keyword)
		if err != nil {
			api.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	})
}
