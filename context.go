package red

import (
	"sync"

	"github.com/valyala/fasthttp"
	stdUrl "net/url"
)

type Ctx struct {
	noCopy noCopy //nolint:unused,structcheck
	*fasthttp.RequestCtx
	next       bool
	pathName   string
	stdUrl     *stdUrl.URL
	deferFuncs []func()
	error      bool
}

var (
	ctxPool sync.Pool
	urlPool sync.Pool
	zeroUrl = &stdUrl.URL{}
)

func AcquireCtx(ctx *fasthttp.RequestCtx) *Ctx {
	v := ctxPool.Get()
	if v == nil {
		redCtx := new(Ctx)
		redCtx.RequestCtx = ctx
		redCtx.stdUrl = AcquireURL(ctx.Request.URI())
		return redCtx
	}

	redCtx := v.(*Ctx)
	redCtx.RequestCtx = ctx
	redCtx.stdUrl = AcquireURL(ctx.Request.URI())

	return redCtx
}

func ReleaseCtx(ctx *Ctx) {
	ctx.next = false

	ctxPool.Put(ctx)
	return
}

func AcquireURL(uri *fasthttp.URI) *stdUrl.URL {
	v := urlPool.Get()
	if v == nil {
		url := new(stdUrl.URL)

		url.Scheme = B2S(uri.Scheme())
		url.Path = B2S(uri.Path())
		url.Host = B2S(uri.Host())
		url.RawQuery = B2S(uri.QueryString())

		return url
	}

	url := v.(*stdUrl.URL)

	url.Scheme = B2S(uri.Scheme())
	url.Path = B2S(uri.Path())
	url.Host = B2S(uri.Host())
	url.RawQuery = B2S(uri.QueryString())
	return url
}

func ReleaseURL(url *stdUrl.URL) {
	*url = *zeroUrl

	urlPool.Put(url)
}

func (ctx *Ctx) Next() {
	ctx.next = true

	return
}

func (ctx *Ctx) IsNext() bool {
	return ctx.next
}

func (ctx *Ctx) PathName() string {
	return ctx.pathName
}

func (ctx *Ctx) URL() *stdUrl.URL {
	return ctx.stdUrl
}

func (ctx *Ctx) Defer(deferFunc func()) {
	ctx.deferFuncs = append(ctx.deferFuncs, deferFunc)
}

func (ctx *Ctx) RequestID() []byte {
	return ctx.Request.Header.Peek(XRequestIDHeader)
}