package goad

import (
	"errors"
	"net/http"
)

type Engine struct {
	Router *Router
}

type Router struct {
	routingTable map[string]func(w http.ResponseWriter, r *http.Request)
}

func (r *Router) Get(path string, handler func(w http.ResponseWriter, r *http.Request)) error {
	if r.routingTable[path] != nil {
		return errors.New("existed")
	}

	r.routingTable[path] = handler
	return nil
}

func (h *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handler := h.Router.routingTable[r.URL.Path]
		if handler == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler(w, r)
		return
	}
}

func New() *Engine {
	return &Engine{
		Router: &Router{
			routingTable: make(map[string]func(w http.ResponseWriter, r *http.Request)),
		},
	}
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}