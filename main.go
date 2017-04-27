package main

import (
	"hecatonhair/SocketEngine"
	"hecatonhair/httpengine"
)

func main() {
	socketEngine := SocketEngine.NewEngine("v1.0")

	httpEngine := httpengine.NewHTTPEngine("v1.0")
	httpEngine.Router.HandlerFunc("GET", "/", socketEngine.Handler.ServeHTTP)
	httpEngine.PowerUp("0.0.0.0", 8181)
}
