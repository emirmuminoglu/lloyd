package red

import (
	"sync"

	"github.com/valyala/fasthttp"
)

type Ctx struct {
	*fasthttp.RequestCtx
	next  bool
	error bool
}

var ctxPool sync.Pool

func AcquireCtx(ctx *fasthttp.RequestCtx) *Ctx {
	v := ctxPool.Get()
	if ctx == nil {
		redCtx := new(Ctx)
		redCtx.RequestCtx = ctx
	}

	redCtx := v.(*Ctx)

	redCtx.RequestCtx = ctx

	return redCtx
}

func ReleaseCtx(ctx *Ctx) {
	ctx.next = false

	ctxPool.Put(ctx)
	return
}

func (ctx *Ctx) Next() {
	ctx.next = true

	return
}

func (ctx *Ctx) IsNext() bool {
	return ctx.next
}
