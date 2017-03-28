package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

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

var srv *http.Server

// Start server with the configuration provided.
func Start(c *config.Config) {
	router := newRouter(c)

	routes, errs := router.generateRoutes()
	if len(errs) != 0 {
		fmt.Printf("%d Error(s) in config: \n", len(errs))
		for i, err := range errs {
			fmt.Printf(" %d: %s\n", i+1, err.Error())
		}
		return
	}

	server := &app{h: routes}

	server.middleware(corsHandler)

	if c.EnableLog {
		server.middleware(logHandler)
	}

	if c.Auth != nil {
		server.middleware(basicAuth(c.Auth.Name, c.Auth.Password))
	}

	if c.JWT != nil {
		server.middleware(jwtVerify(*c.JWT))
	}

	server.middleware(recoverHandler)
	server.middleware(delayHandler(c.Delay))

	if c.Port == 0 {
		c.Port = 8000
	}
	addr := fmt.Sprintf(":%d", c.Port)

	srv = &http.Server{Addr: addr, Handler: server}
	fmt.Println("Starting server at ", addr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

// Stop for stopping s running server
func Stop() {
	if srv == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
