package red

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type RequestHandler func(*Ctx)

func ConvertToFastHTTPHandler(handler RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		rctx := AcquireCtx(ctx)
		defer ReleaseCtx(rctx)

		handler(rctx)

		return
	}

}

type Router struct {
	Group            *fastrouter.Group
	middlewares      []RequestHandler
	deferMiddlewares []RequestHandler
}

func (r *Router) Handle(method, path string, handlers ...RequestHandler) {
	r.Group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		rctx := AcquireCtx(ctx)

		defer func() {
			for _, handler := range r.deferMiddlewares {
				rctx.next = false
				handler(rctx)
				if !rctx.IsNext() {
					break
				}
			}

			ReleaseCtx(rctx)
		}()

		for _, handler := range r.middlewares {
			rctx.next = false
			handler(rctx)
			if !rctx.next {
				return
			}
		}

		for _, handler := range handlers {
			rctx.next = false
			handler(rctx)
			if !rctx.IsNext() {
				return
			}
		}

		return
	})
}

func (r *Router) Use(handlers ...RequestHandler) {
	r.middlewares = append(r.middlewares, handlers...)
}

func (r *Router) Defer(handlers ...RequestHandler) {
	r.deferMiddlewares = append(r.deferMiddlewares, handlers...)
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

func (r *Router) NewGroup(path string) *Router {
	group := r.Group.Group(path)

	return &Router{
		middlewares: r.middlewares,
		Group:       group,
	}
}
