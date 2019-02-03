package static

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func handle() http.Handler {
	stat := New("index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if stat.Response("../htdocs", w, r, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Some-Header", "test")
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}, func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("\n" + `<div>after</div>`))
		}) {
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<div>Error 404!</div>`))
	})
}

func request(t *testing.T, file string) *httptest.ResponseRecorder {
	request, err := http.NewRequest("GET", file, nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handle().ServeHTTP(recorder, request)
	return recorder
}

func TestIndexFile(t *testing.T) {
	r := request(t, "/")
	if s := r.Code; s != http.StatusOK {
		t.Fatalf("handler return wrong status code: got (%v) want (%v)", s, http.StatusOK)
	}
	if c := r.Header().Get("Content-Type"); c != "text/html; charset=utf-8" {
		t.Fatalf("content type header not match: got (%v) want (%v)", c, "text/html; charset=utf-8")
	}
	if r.Body.String() != `Index.html<br />`+"\n"+`<a href="/page.html">page.html</a><br />`+"\n"+`<a href="/sub/index.html">sub/index.html</a>`+"\n\n"+`<div>after</div>` {
		t.Fatalf("bad body response, not match\n%s", r.Body.String())
	}
	if c := r.Header().Get("Some-Header"); c != "test" {
		t.Fatalf("custom header not match: got (%v) want (%v)", c, "test")
	}
}

func TestPageFile(t *testing.T) {
	r := request(t, "/page.html")
	if s := r.Code; s != http.StatusOK {
		t.Fatalf("handler return wrong status code: got (%v) want (%v)", s, http.StatusOK)
	}
	if c := r.Header().Get("Content-Type"); c != "text/html; charset=utf-8" {
		t.Fatalf("content type header not match: got (%v) want (%v)", c, "text/html; charset=utf-8")
	}
	if r.Body.String() != `Page.html<br />`+"\n"+`<a href="/index.html">index.html</a>`+"\n\n"+`<div>after</div>` {
		t.Fatalf("bad body response, not match\n%s", r.Body.String())
	}
	if c := r.Header().Get("Some-Header"); c != "test" {
		t.Fatalf("custom header not match: got (%v) want (%v)", c, "test")
	}
}

func TestIndexFileInSubDir(t *testing.T) {
	r := request(t, "/sub/")
	if s := r.Code; s != http.StatusOK {
		t.Fatalf("handler return wrong status code: got (%v) want (%v)", s, http.StatusOK)
	}
	if c := r.Header().Get("Content-Type"); c != "text/html; charset=utf-8" {
		t.Fatalf("content type header not match: got (%v) want (%v)", c, "text/html; charset=utf-8")
	}
	if r.Body.String() != `Sub index.html<br />`+"\n"+`<a href="/index.html">index.html</a>`+"\n\n"+`<div>after</div>` {
		t.Fatalf("bad body response, not match\n%s", r.Body.String())
	}
	if c := r.Header().Get("Some-Header"); c != "test" {
		t.Fatalf("custom header not match: got (%v) want (%v)", c, "test")
	}
}

func Test404(t *testing.T) {
	r := request(t, "/some.html")
	if s := r.Code; s != http.StatusNotFound {
		t.Fatalf("handler return wrong status code: got (%v) want (%v)", s, http.StatusNotFound)
	}
	if c := r.Header().Get("Content-Type"); c != "text/html; charset=utf-8" {
		t.Fatalf("content type header not match: got (%v) want (%v)", c, "text/html; charset=utf-8")
	}
	if r.Body.String() != `<div>Error 404!</div>` {
		t.Fatalf("bad body response, not match\n%s", r.Body.String())
	}
	if c := r.Header().Get("Some-Header"); c != "" {
		t.Fatalf("custom header not match: got (%v) want (%v)", c, "")
	}
}
