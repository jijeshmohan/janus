package server

import (
	"fmt"
	"net/http"

	"github.com/jijeshmohan/janus/config"
)

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

	if err := http.ListenAndServe(addr, corsHandler(routes)); err != nil {
		fmt.Println(err)
	}
}
