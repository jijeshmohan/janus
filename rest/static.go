package rest

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

type Static struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

func (s *Static) GetEndPoint(rootPath string) (*Endpoint, error) {
	resourcePath := filepath.Join(rootPath, s.Path)

	if !fileExist(resourcePath) {
		return nil, fmt.Errorf("Unable to find static folder %s", resourcePath)
	}

	if !strings.HasPrefix(s.URL, "/") {
		s.URL = "/" + s.URL
	}

	return &Endpoint{URL: s.URL, Method: "GET", Handler: s.getHandle(resourcePath)}, nil
}

func (s Static) getHandle(resourcePath string) http.Handler {
	return http.StripPrefix(s.URL, http.FileServer(http.Dir(resourcePath)))
}
