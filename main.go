package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khushi2411/zocket/database"
	"github.com/khushi2411/zocket/routes"
	
)

func main() {
	// Connect to the database
	database.ConnectDB()
	defer database.CloseDB() // Ensure the database connection is closed on exit

	// Initialize the router
	r := mux.NewRouter()

	// Set up routes
	r.HandleFunc("/create-product", routes.CreateProduct).Methods("POST")
	r.HandleFunc("/get-products", routes.GetAllProducts).Methods("GET")
	r.HandleFunc("/products/{id}", routes.GetProductByID).Methods("GET")

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
