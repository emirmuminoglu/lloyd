package red

import (
	"github.com/valyala/fasthttp"
	adaptor "github.com/valyala/fasthttp/fasthttpadaptor"
	"net/http"
)

func ConvertToFastHTTPHandler(handler RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		rctx := AcquireCtx(ctx)
		defer ReleaseCtx(rctx)

		handler(rctx)

		return
	}

}

func ConvertFastHTTPHandler(handler fasthttp.RequestHandler) RequestHandler {
	return func(c *Ctx) {
		handler(c.RequestCtx)
		return
	}
}

func ConvertStdHTTPHandler(handler http.HandlerFunc) RequestHandler {
	return ConvertFastHTTPHandler(adaptor.NewFastHTTPHandlerFunc(handler))
}
