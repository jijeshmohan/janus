package rest

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTData represents jwt token details.
type JWTData struct {
	URL    string        `json:"url"`
	EXP    int           `json:"exp"`
	Secret string        `json:"secret"`
	Data   jwt.MapClaims `json:"data"`
}

// GetEndPoint send an endpoint for JWT token creation
func (j *JWTData) GetEndPoint(rootPath string) (*Endpoint, error) {
	return &Endpoint{URL: j.URL, Method: "POST", Handler: j.getHandle(rootPath)}, nil
}

func (j *JWTData) getHandle(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := jwt.New(jwt.SigningMethodHS256)
		j.Data["exp"] = time.Now().Add(time.Minute * time.Duration(j.EXP)).Unix()
		token.Claims = j.Data
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(j.Secret))
		if err != nil {
			w.WriteHeader(400)
			return
		}

		w.Header().Set("Content-type", "application/json")
		data := map[string]string{
			"token": tokenString,
		}
		w.WriteHeader(200)
		encoder := json.NewEncoder(w)
		encoder.Encode(data)
	})
}
