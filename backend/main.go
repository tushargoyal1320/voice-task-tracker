/**
 * Voice Task Tracker - Backend Server
 *
 * Main entry point for the Go backend API server
 * Provides RESTful API endpoints for task management (CRUD operations)
 * Includes CORS middleware for cross-origin requests from frontend
 */
package main

import (
	"log"
	"net/http"
	"os"
	"voice-task-tracker/middleware"
	"voice-task-tracker/routes"
)

func main() {
	// Initialize the HTTP router with all API routes
	// Routes are defined in routes/routes.go
	router := routes.SetupRoutes()

	// Apply middleware in order of execution
	// 1. CORSMiddleware: Handles CORS headers for cross-origin requests (frontend at localhost:3000)
	// 2. LoggingMiddleware: Logs all incoming HTTP requests for debugging
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)

	port := os.Getenv("PORT")
	if port == ""{
		port = ":8080"
	}

	// Start the HTTP server on port 8080
	// Server listens for incoming requests and routes them to appropriate handlers
	log.Println("Starting server on port:", port)
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err) // Fatal will exit the program if server fails to start
	}
}
