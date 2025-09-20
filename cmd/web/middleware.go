package main

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

type gzipWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (gw *gzipWriter) Write(b []byte) (int, error) {
	return gw.gzipWriter.Write(b)
}

var gzipPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func (app *application) gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz := gzipPool.Get().(*gzip.Writer)
		defer gzipPool.Put(gz)
		
		gz.Reset(w)
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		
		gw := &gzipWriter{
			ResponseWriter: w,
			gzipWriter:    gz,
		}
		
		next.ServeHTTP(gw, r)
	})
}