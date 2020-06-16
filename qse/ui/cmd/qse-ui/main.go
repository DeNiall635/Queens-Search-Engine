package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	adClient "gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/client"
	searchClient "gitlab2.eeecs.qub.ac.uk/40173800/qse/search/client"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ui/ui"
)

const port = 5000
const searchEndpointEnvName = "SEARCH_ENDPOINT"
const adEndpointEnvName = "AD_ENDPOINT"

func main() {
	// Get provided search endpoint
	searchEndpoint, exists := os.LookupEnv(searchEndpointEnvName)
	if !exists {
		log.Panic("No search endpoint provided")
	}

	// Get provided search endpoint
	adEndpoint, exists := os.LookupEnv(adEndpointEnvName)
	if !exists {
		log.Panic("No ad endpoint provided")
	}

	// Set up search client
	searchClient := searchClient.APISearch{
		Endpoint: searchEndpoint,
	}

	// Set up ad client
	adClient := adClient.APIAdRetrieve{
		Endpoint: adEndpoint,
	}

	// Set up router
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})

	r.Use(cors.Handler)
	// Log requests
	r.Use(middleware.Logger)
	// Monitoring
	r.Use(chiprometheus.NewMiddleware("ui"))
	r.Handle("/metrics", promhttp.Handler())
	// Error handling
	r.NotFound(ui.NotFound())
	r.MethodNotAllowed(ui.MethodNotAllowed())

	r.Get("/", ui.Index())
	r.Get("/search", ui.Search(searchClient, adClient))

	// Set up API server
	srv := http.Server{
		Addr:    fmt.Sprintf(fmt.Sprintf(":%d", port)),
		Handler: r,
	}

	// Start API
	log.Printf("Listening on port %d\n", port)
	srv.ListenAndServe()
}
