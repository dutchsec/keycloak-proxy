// +build go1.8

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

	fmt.Println("1")
	bw := basicWriter{ResponseWriter: w}

	if protoMajor == 2 {
		_, ps := w.(http.Pusher)
		if cn && fl && ps {
			fmt.Println("2")
			return &http2FancyWriter{bw}
		}
	} else {
		_, hj := w.(http.Hijacker)
		_, rf := w.(io.ReaderFrom)
		if cn && fl && hj && rf {
			fmt.Println("3")
			return &httpFancyWriter{bw}
		}
	}
	if fl {
		fmt.Println("4")
		return &flushWriter{bw}
	}

	fmt.Println("5")
	return &bw
}

func (f *http2FancyWriter) Push(target string, opts *http.PushOptions) error {
	return f.basicWriter.ResponseWriter.(http.Pusher).Push(target, opts)
}

var _ http.Pusher = &http2FancyWriter{}
