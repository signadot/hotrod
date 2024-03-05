package frontend

import (
	"embed"
	"net/http"
	"os"
	"testing"
)

//go:embed web_assets/*
var testWebAssetsFS embed.FS

func TestServer(t *testing.T) {
	if os.Getenv("GOTEST_MANUAL") == "" {
		t.Skip()
		return
	}

	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ehlo\n"))
	}))
	m.Handle("/web_assets/", http.FileServer(http.FS(testWebAssetsFS)))

	http.ListenAndServe(":9999", m)

}
