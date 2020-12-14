package red

import (
	"log"
	"net"
	"os"

	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Red struct {
	server *fasthttp.Server
	router *fastrouter.Router
	*Router
}

func New(cfg Config) *Red {

	router := newRouter(cfg)

	if cfg.Name == "" {
		cfg.Name = defaultServerName
	}

	red := &Red{
		server: fasthttpServer(cfg),
		router: router,
	}

	group := router.Group("")

	red.Router = &Router{}

	red.Router.Group = group

	return red
}

func (r *Red) Shutdown(){
	r.server.Shutdown()
}

func newRouter(cfg Config) *fastrouter.Router {
	router := fastrouter.New()

	if cfg.NotFound != nil {
		router.NotFound = ConvertToFastHTTPHandler(cfg.NotFound)
	}

	if cfg.MethodNotAllowed != nil {
		router.MethodNotAllowed = ConvertToFastHTTPHandler(cfg.MethodNotAllowed)
	}

	if cfg.GlobalOPTIONS != nil {
		router.GlobalOPTIONS = ConvertToFastHTTPHandler(cfg.GlobalOPTIONS)
	}

	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stderr, "red: ", 0)
	}

	if cfg.PanicHandler != nil {
		router.PanicHandler = func(ctx *fasthttp.RequestCtx, err interface{}) {
			rctx := AcquireCtx(ctx)
			defer ReleaseCtx(rctx)

			cfg.PanicHandler(rctx, err)
			return
		}
	}

	router.SaveMatchedRoutePath = cfg.SaveMatchedRoutePath

	return router
}

func fasthttpServer(cfg Config) *fasthttp.Server {

	return &fasthttp.Server{
		//		ErrorHandler:                         errHandler,
		Name:                               cfg.Name,
		Concurrency:                        cfg.Concurrency,
		DisableKeepalive:                   cfg.DisableKeepalive,
		ReadBufferSize:                     cfg.ReadBufferSize,
		WriteBufferSize:                    cfg.WriteBufferSize,
		ReadTimeout:                        cfg.ReadTimeout,
		WriteTimeout:                       cfg.WriteTimeout,
		IdleTimeout:                        cfg.IdleTimeout,
		MaxConnsPerIP:                      cfg.MaxConnsPerIP,
		MaxRequestsPerConn:                 cfg.MaxRequestsPerConn,
		MaxKeepaliveDuration:               cfg.MaxKeepaliveDuration,
		TCPKeepalive:                       cfg.TCPKeepalive,
		TCPKeepalivePeriod:                 cfg.TCPKeepalivePeriod,
		MaxRequestBodySize:                 cfg.MaxRequestBodySize,
		ReduceMemoryUsage:                  cfg.ReduceMemoryUsage,
		GetOnly:                            cfg.GetOnly,
		DisablePreParseMultipartForm:       cfg.DisablePreParseMultipartForm,
		LogAllErrors:                       cfg.LogAllErrors,
		DisableHeaderNamesNormalizing:      cfg.DisableHeaderNamesNormalizing,
		SleepWhenConcurrencyLimitsExceeded: cfg.SleepWhenConcurrencyLimitsExceeded,
		NoDefaultServerHeader:              cfg.NoDefaultServerHeader,
		NoDefaultDate:                      cfg.NoDefaultDate,
		NoDefaultContentType:               cfg.NoDefaultContentType,
		ConnState:                          cfg.ConnState,
		Logger:                             cfg.Logger,
		KeepHijackedConns:                  cfg.KeepHijackedConns,
	}
}

func (r *Red) Serve(ln net.Listener) error {
	defer ln.Close()

	r.server.Handler = r.router.Handler

	return r.server.Serve(ln)
}
