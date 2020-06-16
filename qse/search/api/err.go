package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// HTTPError represents when an error has occurred with relevant info, can be marshalled into JSON
type HTTPError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// NotFound handles 404 errors, providing an HTTP error in JSON
func NotFound() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, fmt.Sprintf("Endpoint '%s' not found.", r.URL.Path), http.StatusNotFound)
	})
}

// MethodNotAllowed handles 405 errors, providing an HTTP error in JSON
func MethodNotAllowed() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Error(w, fmt.Sprintf("Method '%s' not allowed for endpoint '%s'.", r.Method, r.URL.Path), http.StatusMethodNotAllowed)
	})
}

// Error constructs, marshals then writes an HTTP error into the response in JSON
func Error(w http.ResponseWriter, message string, code int) {
	json, err := json.Marshal(&HTTPError{
		Message: message,
		Code:    code,
	})
	if err != nil {
		// Should not occur, panic
		log.Panic(err)
		return
	}
	w.WriteHeader(code)
	w.Write(json)
}
