package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error load your configuration:", err.Error())
		return
	}

	port := os.Getenv("PORT")

	fmt.Println("Starting the application...")

	if err := http.ListenAndServe(":"+port, nil); err != nil { // Checked for error
		log.Fatal("The server application is error:", err.Error())
	}
}
