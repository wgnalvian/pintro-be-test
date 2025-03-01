package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"wgnalvian.com/payment-server/controller"
	"wgnalvian.com/payment-server/database"
	"wgnalvian.com/payment-server/routes"
	"wgnalvian.com/payment-server/service"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Load the environment variables
	LoadEnv()

	// Connect to the database
	db := database.ConnectMongo()

	// Load service
	s := &service.UserService{
		Db: db,
	}

	// Load controller
	u := &controller.UserController{
		UserService: s,
	}

	// Initialize the routes
	r := &routes.Route{
		UserController: u,
	}
	router := r.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		database.DisconnectMongo()
		os.Exit(0)
	}()

	router.Run(":3000")
}
