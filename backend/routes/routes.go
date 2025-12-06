/**
 * Routes Package
 *
 * Defines all API routes and maps them to handler functions
 * Uses Gorilla Mux for flexible HTTP routing with URL parameters
 */
package routes

import (
	"voice-task-tracker/handlers"

	"github.com/gorilla/mux"
)

/**
 * SetupRoutes configures all HTTP routes for the application
 *
 * Returns: Configured mux.Router with all API endpoints
 *
 * API Endpoints:
 * - GET    /tasks            - Get all tasks
 * - POST   /tasks            - Create a new task
 * - PUT    /tasks/{id}       - Update a specific task by ID
 * - DELETE /tasks/{id}       - Delete a specific task by ID
 * - POST   /parse-voice      - Parse voice transcript to task fields
 * - GET    /tasks/search     - Search and filter tasks
 */
func SetupRoutes() *mux.Router {
	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Task CRUD endpoints
	// GET: Retrieve all tasks
	router.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")

	// POST: Create a new task (expects JSON body with task fields)
	router.HandleFunc("/tasks", handlers.CreateTask).Methods("POST", "OPTIONS")

	// PUT: Update an existing task by ID (expects JSON body with fields to update)
	router.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT", "OPTIONS")

	// DELETE: Remove a task by ID
	router.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE", "OPTIONS")

	// Voice input endpoint
	// POST: Parse voice transcript into structured task data
	router.HandleFunc("/parse-voice", handlers.ParseVoiceInput).Methods("POST", "OPTIONS")

	// Search endpoint
	// GET: Search and filter tasks by query, status, priority
	router.HandleFunc("/tasks/search", handlers.SearchTasks).Methods("GET")

	return router
}
