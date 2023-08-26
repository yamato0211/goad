package goad

import (
	"log"
	"net/http"
	"strings"
)

type Engine struct {
	Router *Router
}

type Router struct {
	trei Node
}

func (r *Router) Get(path string, handler func(ctx *Context)) error {
	path = strings.TrimSuffix(path, "/")
	existedHandler := r.trei.Search(path)
	if existedHandler != nil {
		panic("this path already used")
	}
	r.trei.Insert(path, handler)
	return nil
}

func (h *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	if r.Method == "GET" {
		path := strings.TrimSuffix(r.URL.Path, "/")
		handler := h.Router.trei.Search(path)
		if handler == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		handler(ctx)
		return
	}
}

func New() *Engine {
	return &Engine{
		Router: &Router{
			trei: NewNode(),
		},
	}
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
