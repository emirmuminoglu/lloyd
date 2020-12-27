package cors

import (
	"github.com/emirmuminoglu/lloyd"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

const strHeaderDelim = ", "

type Config struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	AllowMaxAge      int
	ExposedHeaders   []string
}

func isAllowedOrigin(allowed []string, origin string) bool {
	for _, v := range allowed {
		if v == origin || v == "*" {
			return true
		}
	}

	return false
}

func New(cfg Config) lloyd.RequestHandler {
	allowedHeaders := strings.Join(cfg.AllowedHeaders, strHeaderDelim)
	allowedMethods := strings.Join(cfg.AllowedMethods, strHeaderDelim)
	exposedHeaders := strings.Join(cfg.ExposedHeaders, strHeaderDelim)
	maxAge := strconv.Itoa(cfg.AllowMaxAge)

	return func(ctx *lloyd.Ctx) {
		origin := string(ctx.Request.Header.Peek(fasthttp.HeaderOrigin))

		if !isAllowedOrigin(cfg.AllowedOrigins, origin) {
			ctx.Next()
			return
		}

		ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, origin)

		if cfg.AllowCredentials {
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowCredentials, "true")
		}

		varyHeader := ctx.Response.Header.Peek(fasthttp.HeaderVary)
		if len(varyHeader) > 0 {
			varyHeader = append(varyHeader, strHeaderDelim...)
		}

		varyHeader = append(varyHeader, fasthttp.HeaderOrigin...)
		ctx.Response.Header.SetBytesV(fasthttp.HeaderVary, varyHeader)

		if len(cfg.ExposedHeaders) > 0 {
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlExposeHeaders, exposedHeaders)
		}

		method := string(ctx.Method())
		if method != fasthttp.MethodOptions {
			ctx.Next()
			return
		}

		if len(cfg.AllowedHeaders) > 0 {
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowHeaders, allowedHeaders)
		}

		if len(cfg.AllowedMethods) > 0 {
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowMethods, allowedMethods)
		}

		if cfg.AllowMaxAge > 0 {
			ctx.Response.Header.Set(fasthttp.HeaderAccessControlMaxAge, maxAge)
		}

		ctx.Next()
		return
	}
}
