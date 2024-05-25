package main

import (
	"flag"
	"log"
	"os"

	"crawlquery/html/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	var storagePath string
	flag.StringVar(&storagePath, "storagePath", "/tmp/cd-html", "Path to the storage directory")

	var port string
	flag.StringVar(&port, "port", "8080", "Port to listen on")

	flag.Parse()

	// create storage directory
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	r := gin.Default()

	handler := handler.NewHandler(storagePath)

	r.GET("/pages/:pageID", handler.GetPage)

	r.POST("/pages", handler.StorePage)

	r.Run(":" + port)
}
