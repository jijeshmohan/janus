// Package rest represent different types of rest types and
// its manipulations.
package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Resource respresent a single resource in rest api.
type Resource struct {
	Name    string            `json:"name"`
	Headers map[string]string `json:"headers,omitempty"`
}

// GetEndPoints returns all the endpoints for this resource
// if the specified folder is not present it will send an error
func (r *Resource) GetEndPoints(rootPath string) ([]*Endpoint, error) {
	resourcePath := filepath.Join(rootPath, r.Name)

	if !fileExist(resourcePath) {
		return nil, fmt.Errorf("Unable to find resource folder %s", r.Name)
	}

	endpoints := []*Endpoint{
		&Endpoint{URL: "/" + r.Name, Method: "GET", Handler: r.getHandle(filepath.Join(resourcePath, "index.json"), 200)},
		&Endpoint{URL: "/" + r.Name, Method: "POST", Handler: r.getHandle(filepath.Join(resourcePath, "post.json"), 201)},
		&Endpoint{URL: "/" + r.Name + "/{item}", Method: "GET", Handler: r.getDynamicHandle(resourcePath, 200)},
		&Endpoint{URL: "/" + r.Name + "/{item}", Method: "PUT", Handler: r.getDynamicHandle(resourcePath, 200)},
		&Endpoint{URL: "/" + r.Name + "/{item}", Method: "PATCH", Handler: r.getDynamicHandle(resourcePath, 200)},
		&Endpoint{URL: "/" + r.Name + "/{item}", Method: "DELETE", Handler: r.getDynamicHandle(resourcePath, 200)},
	}
	return endpoints, nil
}

// getHandle provide a static handle for the full filepath
// if the file not present it will send 404
func (r *Resource) getHandle(filePath string, status int) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		data, err := os.Open(filePath)
		if err != nil {
			w.WriteHeader(404)
			return
		}
		defer data.Close()

		w.Header().Set("Content-type", "application/json; charset=utf-8")
		for key, value := range r.Headers {
			if key == "Content-type" {
				continue
			}
			w.Header().Set(key, value)
		}

		w.WriteHeader(status)
		io.Copy(w, data)
	})
}

// getDynamicHandle provide a dynamic handle. it will expect a query
// variable called item and serve the file with the same name
// if the file is not present, it will send 404
func (r *Resource) getDynamicHandle(folderPath string, status int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		item := vars["item"]
		data, err := os.Open(filepath.Join(folderPath, item+".json"))
		if err != nil {
			w.WriteHeader(404)
			return
		}
		defer data.Close()

		w.Header().Set("Content-type", "application/json; charset=utf-8")
		for key, value := range r.Headers {
			if key == "Content-type" {
				continue
			}
			w.Header().Set(key, value)
		}
		w.WriteHeader(status)
		io.Copy(w, data)
	})
}

// check a file or folder exist and return boolean value
func fileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}

	return true
}
