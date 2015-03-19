package server

import "net/http"

// corsHandler middleare to handle cors.
func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
