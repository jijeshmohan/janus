package server

import (
	"errors"
	"sort"

	"github.com/gorilla/mux"
	"github.com/jijeshmohan/janus/config"
	"github.com/jijeshmohan/janus/rest"
)

// router type is for generating routes from the configuration.
type router struct {
	c    *config.Config
	h    *mux.Router
	errs []error
}

// NewRouter create a new router from the configuration provided.
func newRouter(c *config.Config) *router {
	return &router{c: c, h: mux.NewRouter()}
}

// generate routes for all configuration entries.
func (r *router) generateRoutes() (*mux.Router, []error) {

	rootPath := r.c.Path

	//  atleast one resource or url should present.
	if len(r.c.Resources) == 0 && len(r.c.URLs) == 0 {
		r.errs = append(r.errs, errors.New("Please provide atleast one resource or url"))
		return nil, r.errs
	}

	endpoints := make([]*rest.Endpoint, 0, len(r.c.URLs)+(len(r.c.Resources)*7))

	// validate and generate urls URLS
	for _, url := range r.c.URLs {
		u := url
		e, err := u.GetEndPoint(rootPath)
		if err != nil {
			r.errs = append(r.errs, err)
			continue
		}
		endpoints = append(endpoints, e)
	}

	for _, re := range r.c.Resources {
		res := re
		e, err := res.GetEndPoints(rootPath)
		if err != nil {
			r.errs = append(r.errs, err)
			continue
		}
		endpoints = append(endpoints, e...)
	}

	sort.Sort(rest.Endpoints(endpoints))

	for _, en := range endpoints {
		r.h.Handle(en.URL, en.Handler).Methods(en.Method)
	}
	
	return r.h, r.errs
}
