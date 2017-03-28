package rest

import (
	"net/http"
	"strings"
)

// Endpoint defines a single rest endpoint(route) in the server
type Endpoint struct {
	URL     string
	Method  string
	Handler http.Handler
}

// Endpoints represents an Endpoint slice.
type Endpoints []*Endpoint

func (a Endpoints) Len() int {
	return len(a)
}

func (a Endpoints) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Endpoints) Less(i, j int) bool {

	firstDynamic := isDynamic(a[i].URL)
	secondDynamic := isDynamic(a[j].URL)

	if firstDynamic && !secondDynamic {
		return false
	}

	if !firstDynamic && secondDynamic {
		return true
	}

	if len(a[i].URL) != len(a[j].URL) {
		return len(a[i].URL) > len(a[j].URL)
	}

	if a[i].Method != a[j].Method {
		return a[i].Method != "GET"
	}

	if a[i].URL == a[j].URL {
		panic("Two endpoints can't be same")
	}
	return true
}

func isDynamic(url string) bool {
	return strings.Contains(url, "{") && strings.Contains(url, "}")
}
