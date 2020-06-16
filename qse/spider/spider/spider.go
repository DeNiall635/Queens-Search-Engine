package spider

import (
	"database/sql"

	"github.com/gocolly/colly"
)

// Spiderer is the interface for navigating a URL
type Spiderer interface {
	Crawl(host string, url string) error
}

// WebSpider is the spider that crawls web pages and stores the data to a DB
type WebSpider struct {
	DB        *sql.DB
	Collector *colly.Collector
}

// Crawl visits a web page and recursively visits all the links on the page, storing retrieved data to a DB
func (w *WebSpider) Crawl(host string, url string) error {
	var visited []string
	w.Collector.AllowedDomains = []string{host}

	w.Collector.OnHTML("html", func(e *colly.HTMLElement) {
		titleTag := e.DOM.Find("title")
		if titleTag == nil {
			return
		}

		title := titleTag.Text()

		bodyTag := e.DOM.Find("body")
		if bodyTag == nil {
			return
		}

		var exists bool
		err := w.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM page WHERE uri = $1)", e.Request.URL.String()).Scan(&exists)
		if err != nil {
			e.Request.Abort()
		}

		if exists {
			// Already visited
			return
		}

		content := bodyTag.Text()
		_, err = w.DB.Exec("INSERT INTO page (title, uri, index_time, content, tokens) VALUES ($1, $2, now(), $3, to_tsvector($4))",
			title, e.Request.URL.String(), content, content)

		if err != nil {
			e.Request.Abort()
		}
	})

	w.Collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	w.Collector.OnRequest(func(r *colly.Request) {
		target := r.URL.Host + r.URL.Path
		for _, visit := range visited {
			if visit == target {
				r.Abort()
				return
			}
		}
		visited = append(visited, target)
	})

	return w.Collector.Visit(url)
}
