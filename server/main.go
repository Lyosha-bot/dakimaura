package main

import (
	"github.com/joho/godotenv"
	"goserver/server"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	newServer, err := server.NewServer(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	err = newServer.Process()
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
