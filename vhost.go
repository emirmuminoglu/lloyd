package red

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type virtualHost struct {
	Router           *fastrouter.Router
	middlewares      []RequestHandler
	deferMiddlewares []RequestHandler
}

func (r *virtualHost) Handle(method, path string, handlers ...RequestHandler) {
	r.Router.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		rctx := AcquireCtx(ctx)
		rctx.pathName = path

		defer func() {
			for _, handler := range r.deferMiddlewares {
				rctx.next = false
				handler(rctx)
				if !rctx.IsNext() {
					break
				}
			}

			for _, deferFunc := range rctx.deferFuncs {
				deferFunc()
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

func (r *virtualHost) Use(handlers ...RequestHandler) {
	r.middlewares = append(r.middlewares, handlers...)
}

func (r *virtualHost) Defer(handlers ...RequestHandler) {
	r.deferMiddlewares = append(r.deferMiddlewares, handlers...)
}

func (r *virtualHost) GET(path string, handlers ...RequestHandler) {
	r.Handle("GET", path, handlers...)
}

func (r *virtualHost) POST(path string, handlers ...RequestHandler) {
	r.Handle("POST", path, handlers...)
}

func (r *virtualHost) PUT(path string, handlers ...RequestHandler) {
	r.Handle("PUT", path, handlers...)
}

func (r *virtualHost) PATCH(path string, handlers ...RequestHandler) {
	r.Handle("PATCH", path, handlers...)
}

func (r *virtualHost) DELETE(path string, handlers ...RequestHandler) {
	r.Handle("DELETE", path, handlers...)
}

func (r *virtualHost) HEAD(path string, handlers ...RequestHandler) {
	r.Handle("HEAD", path, handlers...)
}

func (r *virtualHost) TRACE(path string, handlers ...RequestHandler) {
	r.Handle("TRACE", path, handlers...)
}

func (r *virtualHost) NewGroup(path string) Router {
	group := r.Router.Group(path)

	return &router{
		middlewares:      r.middlewares,
		deferMiddlewares: r.deferMiddlewares,
		Group:            group,
	}
}

func (r *virtualHost) Handler() fasthttp.RequestHandler {
	return r.Router.Handler
}
