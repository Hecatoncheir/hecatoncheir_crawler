package socket_engine

import (
	"fmt"
	"net/http"
	"golang.org/x/net/websocket"
)

// MessageEvent  is a struct of event for receive from socket server
type MessageEvent struct {
	Message  string
	Data     map[string]interface{}
	ClientID string
}

// NewEngine is a constructor for socket engine
func NewEngine(apiVersion string) *Engine {
	engine := Engine{APIVersion: apiVersion}

	engine.Server = &http.Server{}
	engine.Server.Handler = websocket.Handler(engine.AddConnectedClient)

	engine.Handler = http.HandlerFunc(websocket.Handler(engine.AddConnectedClient).ServeHTTP)

	return &engine
}

// Engine of socket server
type Engine struct {
	APIVersion string
	Server     *http.Server
	Handler    http.HandlerFunc
}

// PowerUp need for start listen port
func (engine *Engine) Listen(host string, port int) {
	engine.Server = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Socket server listen on %v, port:%v \n", host, port)
	engine.Server.ListenAndServe()
}

func (engine *Engine) AddConnectedClient(connection *websocket.Conn) {
	fmt.Print("con")
}
