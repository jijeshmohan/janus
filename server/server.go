package server

import (
	"fmt"
	"net/http"

	"github.com/jijeshmohan/janus/config"
)

type app struct {
	h http.Handler
}

type handler func(http.Handler) http.Handler

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.h.ServeHTTP(w, r)
}

func (a *app) middleware(h handler) {
	a.h = h(a.h)
}

// StartServer starts the server with the configuration provided.
func StartServer(c *config.Config) {

	router := newRouter(c)

	routes, errs := router.generateRoutes()
	if len(errs) != 0 {
		fmt.Printf("%d Error(s) in config: \n", len(errs))
		for i, err := range errs {
			fmt.Printf(" %d: %s\n", i+1, err.Error())
		}
		return
	}

	if c.Port == 0 {
		c.Port = 8000
	}
	addr := fmt.Sprintf(":%d", c.Port)

	fmt.Println("Starting server at ", addr)

	server := &app{h: routes}
	server.middleware(corsHandler)

	if c.EnableLog {
		server.middleware(logHandler)
	}

	if c.Auth != nil {
		server.middleware(basicAuth(c.Auth.Name, c.Auth.Password))
	}

	server.middleware(recoverHandler)
	server.middleware(delayHandler(c.Delay))

	if err := http.ListenAndServe(addr, server); err != nil {
		fmt.Println(err)
	}

}
