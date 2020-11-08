package minirest

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// use the Writer part of gzipResponseWriter to write the output.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// set the HTTP header indicating encoding.
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		fn(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r, p)
	}
}

func writeGzipResp(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(status)
	gz := gzip.NewWriter(w)
	defer gz.Close()
	gz.Write(data)

	return
}
