package lloyd

import (
	"log"
	"net"
	"os"

	fastrouter "github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type Lloyd struct {
	server       *fasthttp.Server
	fastrouter   *fastrouter.Router
	virtualHosts map[string]*virtualHost
	cfg          Config
	*router
}

// New creates a new instance of Lloyd server.
func New(cfg Config) *Lloyd {

	r := newRouter(cfg)

	if cfg.Name == "" {
		cfg.Name = defaultServerName
	}

	l := &Lloyd{
		server:     fasthttpServer(cfg),
		fastrouter: r,
	}

	group := r.Group("")

	l.router = &router{}
	l.cfg = cfg
	l.router.Group = group

	return l
}

// Shutdown shuts the server.
func (l *Lloyd) Shutdown() {
	l.server.Shutdown()
}

// NewVirtualHost creates a virtual host for given hostName.
func (l *Lloyd) NewVirtualHost(hostName string) Router {
	host := new(virtualHost)
	l.virtualHosts[hostName] = host

	host.Router = newRouter(l.cfg)

	return host
}

// Serve serves the server with given listener.
func (l *Lloyd) Serve(ln net.Listener) error {
	defer ln.Close()

	l.server.Handler = func(ctx *fasthttp.RequestCtx) {
		if h := l.virtualHosts[B2S(ctx.URI().Host())]; h != nil {
			h.Handler()(ctx)
		} else {
			l.fastrouter.Handler(ctx)
		}
	}

	return l.server.Serve(ln)
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
