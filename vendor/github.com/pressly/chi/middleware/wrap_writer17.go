// +build go1.7,!go1.8

package middleware

import (
	"fmt"
	"io"
	"net/http"
)

// NewWrapResponseWriter wraps an http.ResponseWriter, returning a proxy that allows you to
// hook into various parts of the response process.
func NewWrapResponseWriter(w http.ResponseWriter, protoMajor int) WrapResponseWriter {
	_, cn := w.(http.CloseNotifier)
	_, fl := w.(http.Flusher)
	fmt.Println("1.1")

	bw := basicWriter{ResponseWriter: w}

	if protoMajor == 2 {
		if cn && fl {
			fmt.Println("1.2")
			return &http2FancyWriter{bw}
		}
	} else {
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		if cn && fl && hj && rf {
			fmt.Println("1.3")
			return &httpFancyWriter{bw}
		}
	}
	if fl {
		fmt.Println("1.4")
		return &flushWriter{bw}
	}

	fmt.Println("1.5")
	return &bw
}
