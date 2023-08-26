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
	treis map[string]*Node
}

func (r *Router) register(method string, path string, handler func(ctx *Context)) error {
	router := r.treis[method]
	path = strings.TrimSuffix(path, "/")
	existedHandler := router.Search(path)
	if existedHandler != nil {
		panic("this path already used")
	}
	router.Insert(path, handler)
	return nil
}

func (r *Router) Get(path string, handler func(ctx *Context)) error {
	return r.register("get", path, handler)
}

func (r *Router) Post(path string, handler func(ctx *Context)) error {
	return r.register("post", path, handler)
}

func (r *Router) Patch(path string, handler func(ctx *Context)) error {
	return r.register("patch", path, handler)
}

func (r *Router) Put(path string, handler func(ctx *Context)) error {
	return r.register("put", path, handler)
}

func (r *Router) Delete(path string, handler func(ctx *Context)) error {
	return r.register("delete", path, handler)
}

func (h *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	router := h.Router.treis[strings.ToLower(r.Method)]
	path := strings.TrimSuffix(r.URL.Path, "/")
	targetNode := router.Search(path)

	if targetNode == nil || targetNode.handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	paramMap := targetNode.ParseParams(r.URL.Path)
	ctx.SetParams(paramMap)
	targetNode.handler(ctx)
}

func New() *Engine {
	return &Engine{
		Router: &Router{
			treis: map[string]*Node{
				"get":    NewNode(),
				"post":   NewNode(),
				"patch":  NewNode(),
				"put":    NewNode(),
				"delete": NewNode(),
			},
		},
	}
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
