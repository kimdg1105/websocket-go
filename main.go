package main

import "websocket-go/network"

func main() {
	n := network.NewNetwork()
	n.StartServer()
}
