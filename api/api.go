package api

import (
	"encoding/json"
	"movies-info-api/omdb"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(apiKey string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/search", handleSearchMovie(apiKey))

	return r
}

func handleSearchMovie(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		if search == "" {
			http.Error(w, "Missing 'search' query parameter", http.StatusBadRequest)
			return
		}

		results, err := omdb.SearchMovies(apiKey, search)
		if err != nil {
			http.Error(w, "Failed to fetch movie data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(results); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
