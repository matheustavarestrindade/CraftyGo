package web_server

import (
	"github.com/gin-gonic/gin"
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

	server.LoadHTMLGlob("templates/*")

	server.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "CraftyGo",
		})
	})

	webServer.instance.Run(":" + port)

	return webServer
}
