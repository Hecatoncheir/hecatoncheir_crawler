package SocketEngine

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// NewEngine is a constructor for socket engine
func NewEngine(apiVersion string) *Engine {
	upgrader := websocket.Upgrader{}

	engine := Engine{APIVersion: apiVersion, headersUpgrader: upgrader, Handler: &ConnectionHandler{}}
	return &engine
}

// Engine of socket server
type Engine struct {
	APIVersion      string
	Server          *http.Server
	headersUpgrader websocket.Upgrader
	Handler         *ConnectionHandler
}

// PowerUp need for start listen port
func (engine *Engine) PowerUp(host string, port int) {
	engine.Server = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Socket server listen on %v, port:%v \n", host, port)

	engine.Server.Handler = engine.Handler
	engine.Server.ListenAndServe()
}

// MessageEvent  is a struct of event for receive from socket server
type MessageEvent struct {
	Message  string
	Data     map[string]interface{}
	ClientID string
}

// ConnectionHandler need for handle new client connection to socket server
type ConnectionHandler struct{}

func (*ConnectionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	fmt.Print(request.Header)
}
