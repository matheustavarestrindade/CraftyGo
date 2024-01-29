package web_server

import (
	"github.com/gin-gonic/gin"
	"github.com/matheustavarestrindade/CraftyGo/frontend"
)

type WebServer struct {
	port     string
	instance *gin.Engine
}

func Start(port string) *WebServer {
	server := gin.Default()

	webServer := &WebServer{
		instance: server,
		port:     port,
	}

	frontend.SvelteKitHandler(server)

	webServer.instance.Run(":" + port)
	return webServer
}
