package main

import (
	"net/http"

	"github.com/vladimirok5959/golang-server-static/static"
)

func main() {
	stat := static.New("index.html")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if stat.Response("./htdocs", w, r, nil, nil) {
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<div>Error 404!</div>`))
	})

	http.ListenAndServe(":8080", nil)
}
