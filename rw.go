package lloyd

import (
	"net/http"
	"sync"
)

var (
	respWriterPool sync.Pool
	zeroRespWriter = &responseWriter{}
)

type responseWriter struct {
	hdr http.Header
	ctx *Ctx
}

func acqRespWriter() *responseWriter {
	v := respWriterPool.Get()
	if v == nil {
		return new(responseWriter)
	}

	return v.(*responseWriter)
}

func relRespWriter(rw *responseWriter) {
	if rw == nil {
		return
	}
	*rw = *zeroRespWriter

	respWriterPool.Put(rw)
}

func (rw *responseWriter) Header() http.Header {
	return rw.hdr
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	return rw.ctx.Write(p)
}

func (rw *responseWriter) StatusCode() int {
	return rw.ctx.rw.StatusCode()
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.ctx.SetStatusCode(statusCode)
}
