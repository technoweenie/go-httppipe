// Package httppipe contains an http.Handler that runs a series of http.Handler
// filters until one of them writes a response.
package httppipe

import (
	"net/http"
)

// A Pipe is a collection of http.Handler filters to run.
type Pipe struct {
	Handlers []http.Handler
	Fallback http.HandlerFunc
}

// New initializes a Pipe with the given handlers.
func New(handlers []http.Handler) *Pipe {
	p := &Pipe{Handlers: handlers}
	p.Fallback = p.serveHTTP
	return p
}

func (p *Pipe) SetFallback(handler http.Handler) {
	p.Fallback = handler.ServeHTTP
}

// ServeHTTP satisfies the http.Handler interface.
func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pipewriter := &pipeWriter{false, w}
	for _, handler := range p.Handlers {
		if handler == nil {
			continue
		}

		handler.ServeHTTP(pipewriter, r)
		if pipewriter.written {
			return
		}
	}

	if !pipewriter.written {
		p.Fallback(w, r)
	}
}

func (p *Pipe) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(":("))
}

type pipeWriter struct {
	written bool
	http.ResponseWriter
}

func (w *pipeWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}
