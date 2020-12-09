package red

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type RequestHandler func(*Ctx)

type Router struct {
	Group       *fastrouter.Group
	middlewares []RequestHandler
}

func (r *Router) Handle(method, path string, handlers ...RequestHandler) {
	r.Group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		rctx := AcquireCtx(ctx)

		for _, handler := range r.middlewares {
			handler(rctx)
			if !rctx.IsNext() {
				break

			}
		}

		for _, handler := range handlers {
			handler(rctx)
			if !rctx.IsNext() {
				break
			}
		}

		ReleaseCtx(rctx)
	})
}

func (r *Router) Use(handlers ...RequestHandler) {
	r.middlewares = append(r.middlewares, handlers...)
}

func (r *Router) GET(path string, handlers ...RequestHandler) {
	r.Handle("GET", path, handlers...)
}

func (r *Router) POST(path string, handlers ...RequestHandler) {
	r.Handle("POST", path, handlers...)
}

func (r *Router) PUT(path string, handlers ...RequestHandler) {
	r.Handle("PUT", path, handlers...)
}

func (r *Router) PATCH(path string, handlers ...RequestHandler) {
	r.Handle("PATCH", path, handlers...)
}

func (r *Router) DELETE(path string, handlers ...RequestHandler) {
	r.Handle("DELETE", path, handlers...)
}

func (r *Router) HEAD(path string, handlers ...RequestHandler) {
	r.Handle("HEAD", path, handlers...)
}

func (r *Router) TRACE(path string, handlers ...RequestHandler) {
	r.Handle("TRACE", path, handlers...)
}
