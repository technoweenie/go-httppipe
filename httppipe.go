package httppipe

import (
	"net/http"
)

type Pipe struct {
	Handlers []http.Handler
	Fallback http.HandlerFunc
}

func New(handlers []http.Handler) *Pipe {
	p := &Pipe{Handlers: handlers}
	p.Fallback = p.serveHTTP
	return p
}

func (p *Pipe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pipewriter := &PipeWriter{false, w}
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

type PipeWriter struct {
	written bool
	http.ResponseWriter
}

func (w *PipeWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}
