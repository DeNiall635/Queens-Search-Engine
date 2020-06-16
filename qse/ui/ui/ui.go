package ui

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/ad"
	adClient "gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/client"
	searchClient "gitlab2.eeecs.qub.ac.uk/40173800/qse/search/client"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/search/search"
)

const templateDir = "templates/"
const queryParam = "q"
const pageParam = "p"

// Page is used with the templates to provide the current page index and available pages
type Page struct {
	Current bool
	Index   int
}

// Index provides a handler for the / endpoint, providing a simple search interface
func Index() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read html
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/index.html", templateDir))
		if err != nil {
			errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}

// Search provides a handler for the /search endpoint, inserting search results and ads
func Search(searchClient searchClient.APISearch, adClient adClient.APIAdRetrieve) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get params
		query := r.URL.Query().Get(queryParam)
		if query == "" {
			redirect(w, "/")
			return
		}

		pageVal := r.URL.Query().Get(pageParam)
		if pageVal == "" {
			redirect(w, fmt.Sprintf("%s&p=0", r.URL.String()))
			return
		}

		page, err := strconv.Atoi(pageVal)
		if err != nil {
			redirect(w, "/")
			return
		}

		// Make search
		results, err := searchClient.Search(query, page)
		if err != nil {
			errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Get ad
		advert, err := adClient.Get(query)
		if err != nil {
			errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Build pages
		var pages []Page
		for i := 0; i < results.PageCount; i++ {
			current := false
			if i == page {
				current = true
			}
			pages = append(pages, Page{
				Current: current,
				Index:   i,
			})
		}

		// Read template
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/search.html", templateDir))
		if err != nil {
			errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Build template
		template, err := template.New("search").Parse(string(data))
		if err != nil {
			errorPage(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)

		// Get search results and ads and insert them
		err = template.Execute(w, struct {
			Query  string
			Result search.Result
			Ad     ad.Advert
			Pages  []Page
		}{
			Query:  query,
			Result: *results,
			Ad:     *advert,
			Pages:  pages,
		})
		if err != nil {
			errorPage(w, http.StatusInternalServerError, err.Error())
		}
	})
}

// NotFound handles 404 page not found errors
func NotFound() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorPage(w, http.StatusNotFound, fmt.Sprintf("Endpoint '%s' not found.", r.URL.Path))
	})
}

// MethodNotAllowed handles 405 method not allowed errors
func MethodNotAllowed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errorPage(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method '%s' not allowed for path '%s'.", r.Method, r.URL.Path))
	})
}

func redirect(w http.ResponseWriter, target string) {
	w.Header().Set("Location", target)
	w.WriteHeader(http.StatusFound)
}

func errorPage(w http.ResponseWriter, code int, message string) {
	// Read template
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/error.html", templateDir))
	if err != nil {
		// Panic, missing error page
		log.Panic(err)
		return
	}

	// Build template
	template, err := template.New("error").Parse(string(data))
	if err != nil {
		// Panic, failed to load error template
		log.Panic(err)
		return
	}

	w.WriteHeader(code)

	err = template.Execute(w, struct {
		Code    int
		Message string
	}{
		code,
		message,
	})
	if err != nil {
		// Panic, failed to load error template
		log.Panic(err)
		return
	}
}
