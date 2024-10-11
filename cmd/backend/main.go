package main

import (
	"encoding/json"
	"log"
	"net/http"
	"flag"
	"fmt"
	"simple_load_balancer/internal/models"
	"simple_load_balancer/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func main() {
	// Define a flag for the port
	port := flag.Int("port", 8080, "port to run the server on")
	flag.Parse()

	var err error
	db, err = database.ConnectMongoDB("mongodb://localhost:27017", "userdb")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/health", handleHealth)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Backend server starting on port %d", *port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = db.Collection("users").InsertOne(r.Context(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "User added"})
	case "GET":
		var user models.User
		err := db.Collection("users").FindOne(r.Context(), map[string]interface{}{}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "No users found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		json.NewEncoder(w).Encode(user)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}