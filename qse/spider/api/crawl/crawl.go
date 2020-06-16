package crawl

import (
	"fmt"
	"net/http"
	"net/url"

	"gitlab2.eeecs.qub.ac.uk/40173800/qse/spider/api"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/spider/spider"
)

const (
	targetParam = "target"
)

// Crawl handles the /crawl endpoint, allows visiting a URL, crawling over links on the page and storing in a DB
func Crawl(crawler spider.Spiderer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := r.URL.Query().Get(targetParam)
		// Check params provided
		if target == "" {
			api.Error(w, "The target URL field 'target' is required.", http.StatusBadRequest)
			return
		}

		uri, err := url.Parse(target)
		if err != nil {
			api.Error(w, fmt.Sprintf("'%s' is not a valid URL", target), http.StatusBadRequest)
		}

		err = crawler.Crawl(uri.Host, target)
		if err != nil {
			api.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
}
