package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gitlab2.eeecs.qub.ac.uk/40173800/qse/search/api"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/search/search"
)

const (
	queryParam = "q"
	pageParam  = "p"
)

// Search handles the /search endpoint, providing simple parameter validation, allows querying results from a DB
func Search(searcher search.Searcher) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get(queryParam)
		// Check params provided
		if query == "" {
			api.Error(w, "The query field 'q' is required.", http.StatusBadRequest)
			return
		}

		pageVal := r.URL.Query().Get(pageParam)
		if pageVal == "" {
			api.Error(w, "The page field 'p' is required.", http.StatusBadRequest)
			return
		}

		// Convert to integers
		page, err := strconv.Atoi(pageVal)
		if err != nil {
			api.Error(w, fmt.Sprintf("The value '%s' is not a valid page number.", pageVal), http.StatusBadRequest)
			return
		}

		// Get search results
		results, err := searcher.Search(query, page)
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
