package socket_engine

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"encoding/json"
	"hecatonhair/crawler"
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
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	engine := Engine{
		APIVersion: apiVersion, Server: &http.Server{}, headersUpgrader: upgrader, Clients: make(map[string]*ConnectedClient)}

	return &engine
}

// Engine of socket server
type Engine struct {
	APIVersion      string
	Server          *http.Server
	Clients         map[string]*ConnectedClient
	clientsMu       sync.Mutex
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

	client := NewConnectedClient(socketConnection)

	engine.clientsMu.Lock()
	engine.Clients[client.ID] = client
	engine.clientsMu.Unlock()

	go engine.listenConnectedClient(client)
}

// listenConnectedClient need for receive and broadcast client messages
func (engine *Engine) listenConnectedClient(client *ConnectedClient) {
	hecatonhair := crawler.NewCrawler()

	go func() {
		for item := range hecatonhair.Items {
			data := map[string]interface{}{"Item": item}

			engine.writeAll("Item from categories of company parsed", data)
		}
	}()

	for event := range client.Channel {
		switch event.Message {
		case "Need api version":

			message := MessageEvent{
				Message: "Version of API",
				Data:    map[string]interface{}{"API version": engine.APIVersion}}

			engine.Clients[event.ClientID].write(message.Message, message.Data)

		case "Get items from categories of company":

			var configuration = crawler.EntityConfig{}
			bytes, _ := json.Marshal(event.Data)
			json.Unmarshal(bytes, &configuration)

			go hecatonhair.RunWithConfiguration(configuration)

		default:
			engine.writeAll(event.Message, event.Data)
		}
	}

	engine.clientsMu.Lock()
	delete(engine.Clients, client.ID)
	engine.clientsMu.Unlock()

}

// writeAll send events to all connected clients
func (engine *Engine) writeAll(message string, details map[string]interface{}) {
	for _, connection := range engine.Clients {
		go connection.write(message, details)
	}
}
