package main

import (
	"dcm-placement-api/placementapi"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MyServer struct{}

func (s *MyServer) GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (s *MyServer) GetHelloName(c *gin.Context, name string) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello " + name})
}

func main() {
	r := gin.Default()
	server := &MyServer{}

	// Register generated routes with Gin
	placementapi.RegisterHandlers(r, server)

	r.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("placement-openapi.yaml")
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
