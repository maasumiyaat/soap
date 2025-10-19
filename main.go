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
	DBPath = "user.db"
	Port   = ":8180"
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
	soapHandler := &handler.UserSOAPHandler{
		UserService: userService,
	}

	// 4. Setup HTTP Server
	http.Handle("/soap/user", soapHandler)

	log.Printf("SOAP Server starting on http://localhost%s/soap/user", Port)
	if err := http.ListenAndServe(Port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
