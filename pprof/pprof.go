package pprof

import (
	"net/http/pprof"
	rtp "runtime/pprof"

	"github.com/emirmuminoglu/lloyd"
)

func Pprof(r lloyd.Router) {
	r.GET("/debug/pprof/cmdline", lloyd.ConvertStdHTTPHandler(pprof.Cmdline))
	r.GET("/debug/pprof/profile", lloyd.ConvertStdHTTPHandler(pprof.Profile))
	r.GET("/debug/pprof/symbol", lloyd.ConvertStdHTTPHandler(pprof.Symbol))
	r.GET("/debug/pprof/trace", lloyd.ConvertStdHTTPHandler(pprof.Trace))
	for _, v := range rtp.Profiles() {
		name := v.Name()
		r.GET("/debug/pprof/"+name, lloyd.ConvertStdHTTPHandler(pprof.Handler(name).ServeHTTP))
	}
}
