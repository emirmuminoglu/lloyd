package pprof

import (
	"net/http/pprof"
	rtp "runtime/pprof"

	"github.com/emirmuminoglu/red"
)

func Pprof(r red.Router) {
	r.GET("/debug/pprof/cmdline", red.ConvertStdHTTPHandler(pprof.Cmdline))
	r.GET("/debug/pprof/profile", red.ConvertStdHTTPHandler(pprof.Profile))
	r.GET("/debug/pprof/symbol", red.ConvertStdHTTPHandler(pprof.Symbol))
	r.GET("/debug/pprof/trace", red.ConvertStdHTTPHandler(pprof.Trace))
	for _, v := range rtp.Profiles() {
		name := v.Name()
		r.GET("/debug/pprof/"+name, red.ConvertStdHTTPHandler(pprof.Handler(name).ServeHTTP))
	}
}
