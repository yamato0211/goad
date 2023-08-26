package goad

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	w      http.ResponseWriter
	r      *http.Request
	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:      w,
		r:      r,
		params: map[string]string{},
	}
}

func (ctx *Context) Json(data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		ctx.w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.w.Header().Set("Content-Type", "application/json")
	ctx.w.WriteHeader(http.StatusOK)
	ctx.w.Write(res)
}

func (ctx *Context) WriteString(data string) {
	ctx.w.WriteHeader(http.StatusOK)
	fmt.Fprint(ctx.w, data)
}

func (ctx *Context) GetQueries() map[string][]string {
	return ctx.r.URL.Query()
}

func (ctx *Context) Query(key string) string {
	values := ctx.GetQueries()
	if target, ok := values[key]; ok {
		if len(target) == 0 {
			return ""
		}
		return target[len(target)-1]
	}
	return ""
}

func (ctx *Context) QueryWithDefault(key string, defaultValue string) string {
	values := ctx.GetQueries()

	if target, ok := values[key]; ok {
		if len(target) == 0 {
			return defaultValue
		}

		return target[len(target)-1]
	}

	return defaultValue
}

func (ctx *Context) SetParams(params map[string]string) {
	ctx.params = params
}

func (ctx *Context) Param(key string) string {
	params := ctx.params

	if v, ok := params[key]; ok {
		return v
	}
	return ""
}
