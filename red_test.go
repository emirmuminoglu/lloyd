package red

import (
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"
)

func Test_New(t *testing.T) {
	notFoundHandler := func(c *Ctx) {}
	methodNotAllowedHandler := func(c *Ctx) {}
	globalOptions := func(c *Ctx) {}
	panicHandler := func(c *Ctx, err interface{}) {}

	cfg := Config{
		NotFound:         notFoundHandler,
		MethodNotAllowed: methodNotAllowedHandler,
		GlobalOPTIONS:    globalOptions,
		PanicHandler:     panicHandler,
	}

	r := New(cfg)

	if r.cfg.NotFound == nil || r.fastrouter.NotFound == nil {
		t.Error("Not found handler is nil")
	}

	if r.cfg.MethodNotAllowed == nil || r.fastrouter.MethodNotAllowed == nil {
		t.Error("Method Not Allowed handler is nil")
	}

	if r.cfg.GlobalOPTIONS == nil || r.fastrouter.GlobalOPTIONS == nil {
		t.Error("Global options handler is nil")
	}

	if r.cfg.PanicHandler == nil || r.fastrouter.PanicHandler == nil {
		t.Error("Panic handler is nil")
	}
}

func TestRed_Serve(t *testing.T) {
	cfg := Config{}
	r := New(cfg)

	ln := fasthttputil.NewInmemoryListener()
	errCh := make(chan error, 1)
	go func() {
		errCh <- r.Serve(ln)
	}()

	time.Sleep(100 * time.Millisecond)

	if err := r.server.Shutdown(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if err := <-errCh; err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if r.server.Handler == nil {
		t.Error("Red.server.Handler is nil")
	}
}

func Benchmark_Handler(b *testing.B) {
	s := New(Config{})

	s.GET("/plaintext", func(ctx *Ctx) { return  })
	s.GET("/json", func(ctx *Ctx) { return  })
	s.GET("/db", func(ctx *Ctx) { return  })
	s.GET("/queries", func(ctx *Ctx) { return  })
	s.GET("/cached-worlds", func(ctx *Ctx) { return  })
	s.GET("/fortunes", func(ctx *Ctx) { return  })
	s.GET("/updates", func(ctx *Ctx) { return  })

	ctx := new(fasthttp.RequestCtx)
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/updates")

	handler := s.fastrouter.Handler

	b.ResetTimer()

	for i := 0; i <= b.N; i++ {
		handler(ctx)
	}
}
