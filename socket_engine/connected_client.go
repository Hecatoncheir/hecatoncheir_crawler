package socket_engine

import (
	//"fmt"
	//"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
	"os"
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

	go func() {
		for {

			inputMessage := MessageEvent{}
			_, messageBytes, err := clientConnection.ReadMessage()

			if err != nil {
				fmt.Fprintf(os.Stdout, "Can't receive message from %s. %v \n", client.ID, err)
				fmt.Fprintf(os.Stdout, "Closed connection of client %s \n", client.ID)
				close(client.Channel)
				break
			}

			json.Unmarshal(messageBytes, &inputMessage)

			inputMessage.ClientID = client.ID
			client.Channel <- inputMessage
		}
	}()

	return &client
}

// write need for send event to client
func (client *ConnectedClient) write(message string, data map[string]interface{}) {
	event := map[string]interface{}{"Message": message, "Data": data}
	websocket.WriteJSON(client.ClientSocket, event)
}
