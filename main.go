package main

import (
	"log"
	"net/http"
	"os"

	"github.com/maasumiyaat/soap/database"
	"github.com/maasumiyaat/soap/handler"
	"github.com/maasumiyaat/soap/model"
	"github.com/maasumiyaat/soap/service"
)

const (
	DBPath   = "user.db"
	HTTPPort = ":8180"
	UDPPort  = ":8181"
)

func main() {
	// 1. Initialize Database
	if err := database.InitDB(DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()
	defer os.Remove(DBPath)

	// 2. Seed Initial Data (Optional, but useful for testing)
	initialUser := &model.User{
		Name:  "Alice Johnson",
		Email: "alice@example.com",
	}
	if err := database.SaveUser(initialUser); err != nil {
		log.Fatalf("Failed to seed initial user: %v", err)
	}
	log.Printf("Seeded user with ID: %d", initialUser.ID)

	// 3. Setup Layers
	userService := &service.UserService{}

	// HTTP SOAP Handler
	httpSoapHandler := &handler.UserSOAPHandler{
		UserService: userService,
	}

	// UDP SOAP Handler
	udpSoapHandler := handler.NewUDPSOAPHandler(userService)

	// 4. Start UDP SOAP Server
	if err := udpSoapHandler.StartUDPServer("localhost" + UDPPort); err != nil {
		log.Fatalf("Failed to start UDP SOAP server: %v", err)
	}
	defer udpSoapHandler.Stop()

	// 5. Setup HTTP Server
	http.Handle("/soap/user", httpSoapHandler)

	log.Printf("HTTP SOAP Server starting on http://localhost%s/soap/user", HTTPPort)
	log.Printf("UDP SOAP Server listening on localhost%s", UDPPort)

	if err := http.ListenAndServe(HTTPPort, nil); err != nil {
		log.Fatalf("HTTP Server failed: %v", err)
	}
}
