package main

import (
	"net/http"
	"log"
	
	"example.com/server/users"
)

func main() {
	http.HandleFunc("/api/users", users.GetUsersHandler)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}




