package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engine *gin.Engine
}

func NewNetwork() *Network {
	n := &Network{
		engine: gin.New(),
	}

	n.engine.Use(gin.Logger())
	n.engine.Use(gin.Recovery())
	n.engine.Use(cors.New(cors.Config{
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"*"},
		AllowWebSockets: true,
	}))

	return n
}

func (n *Network) StartServer() error {
	log.Println("StartServer ...")
	return n.engine.Run(":8080")
}
