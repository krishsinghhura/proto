package main

import (
	"fmt"
	"log"
	"verification/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
func main() {

	fmt.Println("Starting server...")
	r := gin.Default()
	r.POST("/upload", handlers.UploadHandler)

	r.POST("/verify", handlers.VerifyHandler)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
