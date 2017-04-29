package main

import (
	"hecatonhair/http_engine"
	socket "hecatonhair/socket_engine"
)

func main() {
	socketEngine := socket.NewEngine("v1.0")

	httpEngine := http_engine.NewHTTPEngine("v1.0")
	httpEngine.Router.HandlerFunc("GET", "/", socketEngine.AddConnectedClient)
	httpEngine.PowerUp("0.0.0.0", 8181)
}
