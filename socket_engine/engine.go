package socket_engine

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"encoding/json"
)

// MessageEvent  is a struct of event for receive from socket server
type MessageEvent struct {
	Message  string
	Data     map[string]interface{}
	ClientID string
}

// NewEngine is a constructor for socket engine
func NewEngine(apiVersion string) *Engine {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	engine := Engine{
		APIVersion: apiVersion, Server: &http.Server{}, headersUpgrader: upgrader}

	return &engine
}

// Engine of socket server
type Engine struct {
	APIVersion      string
	Server          *http.Server
	headersUpgrader websocket.Upgrader
}

// PowerUp need for start listen port
func (engine *Engine) Listen(host string, port int) {
	engine.Server = &http.Server{Addr: fmt.Sprintf("%v:%v", host, port)}
	fmt.Printf("Socket server listen on %v, port:%v \n", host, port)
	engine.Server.Handler = http.HandlerFunc(engine.AddConnectedClient)
	engine.Server.ListenAndServe()
}

// AddConnectedClient is a handler for new socket connections
func (engine *Engine) AddConnectedClient(response http.ResponseWriter, request *http.Request) {

	socketConnection, err := engine.headersUpgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, messageBytes, err := socketConnection.ReadMessage()
		if err != nil {
			return
		}

		event := &MessageEvent{}
		json.Unmarshal(messageBytes, event)

		socketConnection.WriteJSON(event)
	}
}
