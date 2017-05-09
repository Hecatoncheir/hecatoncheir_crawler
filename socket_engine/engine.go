package socket_engine

import (
	"encoding/json"
	"fmt"
	"hecatoncheir/crawler"
	"hecatoncheir/crawler/mvideo"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// MessageEvent  is a struct of event for receive from socket server
type MessageEvent struct {
	Message  string
	Data     interface{}
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
		APIVersion:      apiVersion,
		Server:          &http.Server{},
		headersUpgrader: upgrader,
		Clients:         make(map[string]*ConnectedClient)}

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

	for event := range client.Channel {
		switch event.Message {
		case "Need api version":

			message := MessageEvent{
				Message: "Version of API",
				Data:    map[string]interface{}{"API version": engine.APIVersion}}

			engine.Clients[event.ClientID].Write(message.Message, message.Data)

		case "Get items from categories of company":

			var company = crawler.Company{}
			dataBytes, _ := json.Marshal(event.Data)
			json.Unmarshal(dataBytes, &company)

			if company.Iri == "http://www.mvideo.ru/" {
				hecatonhair := mvideo.NewCrawler()

				go func() {
					for item := range hecatonhair.Items {
						data := map[string]interface{}{"Item": item}

						engine.WriteAll("Item from categories of company parsed", data)
					}
				}()

				var configuration = mvideo.EntityConfig{}
				json.Unmarshal(dataBytes, &configuration)

				go hecatonhair.RunWithConfiguration(configuration)
			}

		default:
			engine.WriteAll(event.Message, event.Data)
		}
	}

	engine.clientsMu.Lock()
	delete(engine.Clients, client.ID)
	engine.clientsMu.Unlock()

}

// WriteAll send events to all connected clients
func (engine *Engine) WriteAll(message string, data interface{}) {
	for _, connection := range engine.Clients {
		go connection.Write(message, data)
	}
}
