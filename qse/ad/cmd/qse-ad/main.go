package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq" // Postgres DB driver
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/ad"
	"gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/api"
	adapi "gitlab2.eeecs.qub.ac.uk/40173800/qse/ad/api/ad"
)

const port = 5000
const connStringEnvName = "AD_DB"

func main() {
	// Get provided connection string for DB
	connString, exists := os.LookupEnv(connStringEnvName)
	if !exists {
		log.Panic("No ad DB connection string provided")
	}

	// Connect to DB
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// Set up advert server
	server := ad.DBAdServer{
		DB: db,
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
	r.Use(chiprometheus.NewMiddleware("ad"))
	r.Handle("/metrics", promhttp.Handler())
	// Error handling
	r.NotFound(api.NotFound())
	r.MethodNotAllowed(api.MethodNotAllowed())

	r.Get("/ad", adapi.Ad(&server))
	r.Post("/ad", adapi.PostAd(&server))
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Health chceck
		w.WriteHeader(http.StatusOK)
	}))

	// Set up API server
	srv := http.Server{
		Addr:    fmt.Sprintf(fmt.Sprintf(":%d", port)),
		Handler: r,
	}

	// Start API
	log.Printf("Listening on port %d\n", port)
	srv.ListenAndServe()
}
