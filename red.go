package red

import (
	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Red struct {
	server *fasthttp.Server
	router *fastrouter.Router
	*Router
}

func New(cfg Config) *Red {
	router := fastrouter.New()

	red := &Red{
		server: fasthttpServer(cfg),
		router: router,
	}

	group := router.Group("/")

	red.Router.Group = group

	return red
}

func (r *Red) NewGroup(path string) {
	r.router.Group(path)
}

func fasthttpServer(cfg Config) *fasthttp.Server {
	return &fasthttp.Server{
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
		KeepHijackedConns:                  cfg.KeepHijackedConns,
	}
}
