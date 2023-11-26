package proxy

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/telanflow/mps"
)

var (
	upgrader     = websocket.Upgrader{}
	endPointAddr = "ws.okex.com:8443"
)

// run a endPoint websocket server
func runWebsocketServer() {
	http.ListenAndServe(endPointAddr, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}))
}

// A simple proxy websocket server
func Test_okx_proxy(t *testing.T) {
	// quit signal
	quitSignChan := make(chan os.Signal)
	signal.Notify(quitSignChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)

	// start endPoint websocket server
	go runWebsocketServer()

	// start proxy websocket server
	websocketHandler := mps.NewWebsocketHandler()
	websocketHandler.Transport().Proxy = func(request *http.Request) (*url.URL, error) {
		// endPoint websocket server
		return url.Parse("ws://" + endPointAddr)
	}
	srv := &http.Server{
		Addr:    "54.169.108.72:8090",
		Handler: websocketHandler,
	}
	go func() {
		log.Printf("WebsocketProxy started listen: ws://%s", srv.Addr)
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		if err != nil {
			quitSignChan <- syscall.SIGKILL
			log.Fatalf("WebsocketProxy start fail: %v", err)
		}
	}()

	<-quitSignChan
	_ = srv.Close()
	log.Fatal("WebsocketProxy server stop!")
}
