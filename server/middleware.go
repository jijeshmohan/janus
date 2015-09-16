package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jijeshmohan/janus/rest"
)

// corsHandler middleare to handle cors.
func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// delayHandler middleare to handle delay in response.
func delayHandler(delay int) handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if delay > 0 {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
			h.ServeHTTP(w, r)
		})
	}
}

// basicAuth middleware for handling basic auth request.
func basicAuth(username string, password string) handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			name, pass, ok := req.BasicAuth()
			if !ok || !(username == name && password == pass) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			h.ServeHTTP(w, req)
		})
	}
}

// recoverHandler middleware for panic recovery
func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.Proto, r.URL)
		h.ServeHTTP(w, r)
	})
}

// jwtVerify middleware for handling jwt verification in request.
func jwtVerify(j rest.JWTData) handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.URL.RequestURI() == j.URL {
				h.ServeHTTP(w, req)
				return
			}
			auth := req.Header.Get("Authorization")
			if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				auth = auth[7:]
			}
			token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
				return []byte(j.Secret), nil
			})
			if err == nil && token.Valid {
				h.ServeHTTP(w, req)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		})
	}
}
