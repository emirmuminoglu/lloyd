package red

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type RequestHandler func(*Ctx)

type Router interface {
	// Handle registers given request handlers with the given path and method
	// There are shortcuts for some methods you can use them.
	Handle(method string, path string, handlers ...RequestHandler)

	// Use registers given middleware handlers to router
	// Given handlers will be executed by given order
	Use(handlers ...RequestHandler)

	// Defer registers given handlers to router
	// Given handlers will be executed after middlewares and handlers by given order
	Defer(handlers ...RequestHandler)

	// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handlers)
	GET(path string, handlers ...RequestHandler)

	// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handlers)
	POST(path string, handlers ...RequestHandler)

	// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handlers)
	PUT(path string, handlers ...RequestHandler)

	// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handlers)
	PATCH(path string, handlers ...RequestHandler)

	// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handlers)
	DELETE(path string, handlers ...RequestHandler)

	// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handlers)
	HEAD(path string, handlers ...RequestHandler)

	// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handlers)
	TRACE(path string, handlers ...RequestHandler)

	// Handler gives routers request handler.
	// Returns un-nil function only if the router is a virtual host router.
	Handler() fasthttp.RequestHandler

	// NewGroup creates a subrouter for given path.
	NewGroup(path string) Router
}

type router struct {
	Group            *fastrouter.Group
	path             string
	middlewares      []RequestHandler
	deferMiddlewares []RequestHandler
}

// Handle registers given request handlers with the given path and method
// There are shortcuts for some methods you can use them.
func (r *router) Handle(method, path string, handlers ...RequestHandler) {
	r.Group.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		rctx := AcquireCtx(ctx)
		rctx.pathName = r.path + path

		defer func() {
			for _, handler := range r.deferMiddlewares {
				rctx.next = false
				handler(rctx)
				if !rctx.next {
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
			if !rctx.next {
				return
			}
		}

		return
	})
}

// Use registers given middleware handlers to router
// Given handlers will be executed by given order
func (r *router) Use(handlers ...RequestHandler) {
	r.middlewares = append(r.middlewares, handlers...)
}

// Defer registers given handlers to router
// Given handlers will be executed after middlewares and handlers by given order
func (r *router) Defer(handlers ...RequestHandler) {
	r.deferMiddlewares = append(r.deferMiddlewares, handlers...)
}

// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handlers)
func (r *router) GET(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodGet, path, handlers...)
}

// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handlers)
func (r *router) POST(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodPost, path, handlers...)
}

// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handlers)
func (r *router) PUT(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodPut, path, handlers...)
}

// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handlers)
func (r *router) PATCH(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodPatch, path, handlers...)
}

// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handlers)
func (r *router) DELETE(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodDelete, path, handlers...)
}

// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handlers)
func (r *router) HEAD(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodHead, path, handlers...)
}

// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handlers)
func (r *router) TRACE(path string, handlers ...RequestHandler) {
	r.Handle(fasthttp.MethodTrace, path, handlers...)
}

// Handler gives routers request handler.
// Returns un-nil function only if the router is a virtual host router.
func (r *router) Handler() fasthttp.RequestHandler {
	return nil
}

// NewGroup creates a subrouter for given path.
func (r *router) NewGroup(path string) Router {
	group := r.Group.Group(path)

	return &router{
		middlewares:      r.middlewares,
		deferMiddlewares: r.deferMiddlewares,
		path:             r.path + path,
		Group:            group,
	}
}
