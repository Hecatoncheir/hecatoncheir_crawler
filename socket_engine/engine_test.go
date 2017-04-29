package socket_engine

import (
	"sync"
	"testing"

	"golang.org/x/net/websocket"
)

var once sync.Once

func startUpSocketServer() {
	engine := NewEngine("v1.0")
	engine.Listen("localhost", 8181)
	defer engine.Server.Close()
}

func TestSocketServerCanHandleEvents(test *testing.T) {
	go once.Do(startUpSocketServer)

	socketConnection, err := websocket.Dial("ws://localhost:8181", "", "http://localhost:8181")
	if err != nil {
		test.Error(err)
	}

	inputMessage := make(chan MessageEvent)

	go func() {
		defer socketConnection.Close()
		defer close(inputMessage)

		for {
			messageFromServer := MessageEvent{}
			err = websocket.JSON.Receive(socketConnection, &messageFromServer)

			if err != nil {
				test.Error(err)
				break
			}

			inputMessage <- messageFromServer
		}
	}()

	messageToServer := MessageEvent{Message: "Need api version"}
	err = websocket.JSON.Send(socketConnection, messageToServer)

	if err != nil {
		test.Error(err)
	}

	for messageFromServer := range inputMessage {
		if messageFromServer.Message != "Version of API" ||
				messageFromServer.Data["API version"] != "v1.0" {
			test.Fail()
		}
		break
	}
}
