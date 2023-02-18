package pkg

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGitlab(t *testing.T) {
	raw, err := os.ReadFile("../testdata/event_push.json")
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/gitlab", bytes.NewBuffer(raw))
	w := httptest.NewRecorder()
	engine := CreateServer(ServerConfig{
		Port:     9448,
		SibylUrl: "http://127.0.0.1:9876",
	})
	engine.ServeHTTP(w, req)

	if w.Code != 200 {
		panic(nil)
	}
	time.Sleep(5 * time.Second)
}

func TestSendTestEvent(t *testing.T) {
	t.Skip()
	raw, err := os.ReadFile("../testdata/event_push.json")
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:9448/api/v1/gitlab", bytes.NewBuffer(raw))
	_, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
}
