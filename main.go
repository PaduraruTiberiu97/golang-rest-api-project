package main

import (
	"apiproject/db"
	"apiproject/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	server := gin.Default()
	routes.RegisterRoutes(server)

	if err := server.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
