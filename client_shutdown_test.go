package remotedialer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestClientCancellationDuringStartup(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()
	for i := 0; i < 100; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := ConnectToProxy(ctx, "ws"+strings.TrimPrefix(server.URL, "http"), nil,
			func(string, string) bool { return false }, nil,
			func(context.Context, *Session) error { cancel(); return nil })
		cancel()
		if err != nil {
			t.Fatal(err)
		}
	}
}
