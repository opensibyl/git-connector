package pkg

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGitlab(t *testing.T) {
	raw, err := os.ReadFile("../testdata/event_push.json")
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/gitlab", bytes.NewBuffer(raw))
	w := httptest.NewRecorder()
	engine := CreateServer()
	engine.ServeHTTP(w, req)

	if w.Code != 200 {
		panic(nil)
	}
}
