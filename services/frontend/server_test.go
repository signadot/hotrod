package frontend

import (
	"embed"
	"net/http"
	"testing"
)

//go:embed web_assets/*
var testWebAssetsFS embed.FS

//go:embed templates/*
var testTplFS embed.FS

func TestServer(t *testing.T) {
	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ehlo\n"))
	}))
	m.Handle("/web_assets/", http.FileServer(http.FS(testWebAssetsFS)))

	http.ListenAndServe(":9999", m)

}
