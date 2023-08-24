package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ankeshnirala/go/authentication/api"
	"github.com/ankeshnirala/go/authentication/storage"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "auth-service ", log.LstdFlags)

	err := godotenv.Load("app.env")
	if err != nil {
		logger.Println(err)
	}

	appPort := os.Getenv("APP_PORT")

	listenAddr := flag.String("listenaddr", appPort, "the server address")
	flag.Parse()

	mongoStore, err := storage.NewMongoStore()
	if err != nil {
		logger.Println(err)
	}

	server := api.NewServer(logger, *listenAddr, mongoStore)
	msg := fmt.Sprintf("Server is running on port %s", *listenAddr)
	logger.Println(msg)

	logger.Fatal(server.Start())
}
