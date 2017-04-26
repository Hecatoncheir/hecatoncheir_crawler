package hecatonhair

import "hecatonhair/httpengine"

func main() {
	httpServer := httpengine.NewHTTPEngine("v1.0")
	httpServer.PowerUp("0.0.0.0", 8181)
}