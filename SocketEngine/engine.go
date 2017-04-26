package SocketEngine

import (
	"github.com/gorilla/websocket"
	"fmt"
	"net/http"
)

// NewEngine is a constructor for socket engine
func NewEngine(apiVersion string) *Engine {
	upgrader := websocket.Upgrader{}

	engine := Engine{APIVersion: apiVersion, headersUpgrader: upgrader, handler: &ConnectionHandler{}}
	return &engine
}

// Engine of socket server
type Engine struct {
	APIVersion      string
	Server          *http.Server
	headersUpgrader websocket.Upgrader
	handler         *ConnectionHandler
}

// PowerUp need for start listen port
func (engine *Engine) PowerUp(host string, port int) {
	engine.Server = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Socket server listen on %v, port:%v \n", host, port)

	engine.Server.Handler = engine.handler
	engine.Server.ListenAndServe()
}

// ConnectionHandler need for handle new client conenction to socket server
type ConnectionHandler struct{}

func (*ConnectionHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	fmt.Print(request.Header)
}
