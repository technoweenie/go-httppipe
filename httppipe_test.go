package httppipe

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPipe(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	p := New([]http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("A", "1")
		}),
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("B", "2")
		}),
		setup.Handler(),
	})

	res := setup.Call(p)
	if res.StatusCode != 200 {
		t.Errorf("Bad response status: %d", res.StatusCode)
	}

	if v := res.Header.Get("A"); v != "1" {
		t.Errorf("Bad A value: %s", v)
	}

	if v := res.Header.Get("B"); v != "2" {
		t.Errorf("Bad B value: %s", v)
	}

	by, _ := ioutil.ReadAll(res.Body)
	if string(by) != "OK" {
		t.Errorf("Bad response body: %s", string(by))
	}
}

func TestSkipNilPipe(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	p := New([]http.Handler{
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("A", "1")
		}),
		nil,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("B", "2")
		}),
		setup.Handler(),
	})

	res := setup.Call(p)
	if res.StatusCode != 200 {
		t.Errorf("Bad response status: %d", res.StatusCode)
	}

	if v := res.Header.Get("A"); v != "1" {
		t.Errorf("Bad A value: %s", v)
	}

	if v := res.Header.Get("B"); v != "2" {
		t.Errorf("Bad B value: %s", v)
	}

	by, _ := ioutil.ReadAll(res.Body)
	if string(by) != "OK" {
		t.Errorf("Bad response body: %s", string(by))
	}
}

type PipeSetup struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	t      *testing.T
}

func Setup(t *testing.T) *PipeSetup {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	return &PipeSetup{srv, mux, t}
}

func (s *PipeSetup) Teardown() {
	s.Server.Close()
}

func (s *PipeSetup) Call(p *Pipe) *http.Response {
	s.Mux.Handle("/", p)
	res, err := http.Get(s.Server.URL + "/")
	if err != nil {
		s.t.Fatalf("Error making request: %s", err)
	}
	return res
}

func (s *PipeSetup) Handler() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	return http.HandlerFunc(f)
}
