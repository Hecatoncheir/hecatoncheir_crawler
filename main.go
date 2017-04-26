package main

import "hecatonhair/httpengine"
import "hecatonhair/SocketEngine"

func main() {
	httpServer := httpengine.NewHTTPEngine("v1.0")
	httpServer.PowerUp("0.0.0.0", 8181)

	socketEngine := SocketEngine.NewEngine("v1.0")
	socketEngine.PowerUp("0.0.0.0", 8182)
}