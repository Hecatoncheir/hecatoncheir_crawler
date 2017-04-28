package socket_engine

import (
	//"fmt"
	//"os"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

// ConnectedClient of socket connection
type ConnectedClient struct {
	ID           string
	Channel      chan MessageEvent
	ClientSocket *websocket.Conn
}

// NewConnectedClient for constructor for ConnectedClient
func NewConnectedClient(clientConnection *websocket.Conn) *ConnectedClient {

	clientID, _ := uuid.NewUUID()
	client := ConnectedClient{ID: clientID.String(), ClientSocket: clientConnection, Channel: make(chan MessageEvent)}

	//go func() {
	//	for {
	//		defer close(client.Channel)
	//
	//		inputMessage := Event{}
	//		err := websocket.JSON.Receive(clientConnection, &inputMessage)
	//
	//		if err != nil {
	//			fmt.Fprintf(os.Stdout, "Can't receive message from %s. %v", client.ID, err)
	//			break
	//		}
	//
	//		inputMessage.ClientID = client.ID
	//		client.Channel <- inputMessage
	//	}
	//}()

	return &client
}

func (client *ConnectedClient) write(message string, data map[string]interface{}) {
	event := MessageEvent{Message: message, Data: data}
	websocket.JSON.Send(client.ClientSocket, event)
}
