package lloyd

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

type Config struct {
	//FastHTTP Settings
	ErrorHandler                       func(*Ctx, error)
	Name                               string
	Concurrency                        int
	DisableKeepalive                   bool
	ReadBufferSize                     int
	WriteBufferSize                    int
	ReadTimeout                        time.Duration
	WriteTimeout                       time.Duration
	IdleTimeout                        time.Duration
	MaxConnsPerIP                      int
	MaxRequestsPerConn                 int
	MaxKeepaliveDuration               time.Duration
	TCPKeepalive                       bool
	TCPKeepalivePeriod                 time.Duration
	MaxRequestBodySize                 int
	ReduceMemoryUsage                  bool
	GetOnly                            bool
	DisablePreParseMultipartForm       bool
	LogAllErrors                       bool
	DisableHeaderNamesNormalizing      bool
	SleepWhenConcurrencyLimitsExceeded time.Duration
	NoDefaultServerHeader              bool
	NoDefaultDate                      bool
	NoDefaultContentType               bool
	ConnState                          func(net.Conn, fasthttp.ConnState)
	Logger                             fasthttp.Logger
	KeepHijackedConns                  bool

	//Router Settigs
	SaveMatchedRoutePath bool
	GlobalOPTIONS        RequestHandler
	NotFound             RequestHandler
	MethodNotAllowed     RequestHandler
	PanicHandler         func(*Ctx, interface{})
}
