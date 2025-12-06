/**
 * Task Handlers Package
 *
 * Contains all HTTP handler functions for task CRUD operations
 * Manages in-memory task storage with thread-safe operations
 */
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// Task represents a single task item with full details
// JSON tags specify how the struct is serialized to/from JSON
type Task struct {
	ID          int    `json:"id"`                    // Unique identifier (auto-incremented)
	Title       string `json:"title"`                 // Task title/summary
	Description string `json:"description,omitempty"` // Detailed task description (optional)
	Status      string `json:"status"`                // Task status: "To Do", "In Progress", "Done"
	Priority    string `json:"priority"`              // Priority level: "Low", "Medium", "High", "Urgent"
	DueDate     string `json:"dueDate,omitempty"`     // Due date in ISO format (optional)
	CreatedAt   string `json:"createdAt"`             // Creation timestamp
	UpdatedAt   string `json:"updatedAt"`             // Last update timestamp
}

// VoiceParseRequest represents the request body for voice-to-task parsing
type VoiceParseRequest struct {
	Transcript string `json:"transcript"` // Raw speech-to-text transcript
}

// VoiceParseResponse represents the AI-parsed task data
type VoiceParseResponse struct {
	Transcript  string `json:"transcript"`            // Original transcript
	Title       string `json:"title"`                 // Extracted task title
	Description string `json:"description,omitempty"` // Extracted description
	Priority    string `json:"priority"`              // Extracted priority
	DueDate     string `json:"dueDate,omitempty"`     // Parsed due date
	Status      string `json:"status"`                // Default status
}

// Global variables for in-memory task storage
var (
	tasks  = []Task{} // Slice to store all tasks
	nextID = 1        // Counter for generating unique task IDs
	mu     sync.Mutex // Mutex to protect concurrent access to tasks slice
)

/**
 * GetTasks returns all tasks from memory
 *
 * HTTP Method: GET
 * Endpoint: /tasks
 * Response: JSON array of all tasks
 *
 * Thread-safe: Uses mutex lock to prevent race conditions
 */
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Lock mutex to ensure thread-safe read
	mu.Lock()
	defer mu.Unlock() // Unlock after function returns

	// Encode tasks slice as JSON and write to response
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/**
 * CreateTask creates a new task and adds it to memory
 *
 * HTTP Method: POST
 * Endpoint: /tasks
 * Request Body: JSON object with task fields (title, description, status, priority, dueDate)
 * Response: JSON object of newly created task (with auto-generated ID and timestamps)
 * Status: 201 Created on success, 400 Bad Request on error
 *
 * Thread-safe: Uses mutex lock to prevent race conditions during write
 */
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Create request")
	// Decode JSON request body into Task struct
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set defaults for required fields
	if task.Status == "" {
		task.Status = "To Do"
	}
	if task.Priority == "" {
		task.Priority = "Medium"
	}

	// Lock mutex for thread-safe write operation
	mu.Lock()
	task.ID = nextID // Assign unique ID
	nextID++         // Increment ID counter for next task

	// Set timestamps
	timestamp := getCurrentTimestamp()
	task.CreatedAt = timestamp
	task.UpdatedAt = timestamp

	tasks = append(tasks, task) // Add task to slice
	mu.Unlock()

	// Return 201 Created status with the new task
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

/**
 * UpdateTask updates an existing task's properties
 *
 * HTTP Method: PUT
 * Endpoint: /tasks/{id}
 * URL Parameter: id (task ID)
 * Request Body: JSON object with fields to update (completed, title)
 * Response: JSON object of updated task
 * Status: 200 OK on success, 400 Bad Request for invalid ID, 404 Not Found if task doesn't exist
 *
 * Thread-safe: Uses mutex lock to prevent race conditions during update
 */
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract task ID from URL path parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) // Convert string ID to integer
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Decode request body with updated task data
	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Lock mutex for thread-safe read and write
	mu.Lock()
	defer mu.Unlock()

	// Find task by ID and update its properties
	for i, task := range tasks {
		if task.ID == id {
			// Update all fields if provided
			if updatedTask.Title != "" {
				tasks[i].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {
				tasks[i].Description = updatedTask.Description
			}
			if updatedTask.Status != "" {
				tasks[i].Status = updatedTask.Status
			}
			if updatedTask.Priority != "" {
				tasks[i].Priority = updatedTask.Priority
			}
			if updatedTask.DueDate != "" {
				tasks[i].DueDate = updatedTask.DueDate
			}

			// Update timestamp
			tasks[i].UpdatedAt = getCurrentTimestamp()

			// Return updated task as JSON
			if err := json.NewEncoder(w).Encode(tasks[i]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// Task with given ID not found
	http.Error(w, "Task not found", http.StatusNotFound)
}

/**
 * DeleteTask removes a task from memory
 *
 * HTTP Method: DELETE
 * Endpoint: /tasks/{id}
 * URL Parameter: id (task ID)
 * Response: No content (empty response body)
 * Status: 204 No Content on success, 400 Bad Request for invalid ID, 404 Not Found if task doesn't exist
 *
 * Thread-safe: Uses mutex lock to prevent race conditions during delete
 */
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract task ID from URL path parameter
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) // Convert string ID to integer
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Lock mutex for thread-safe read and write
	mu.Lock()
	defer mu.Unlock()

	// Find task by ID and remove it from slice
	for i, task := range tasks {
		if task.ID == id {
			// Remove task by concatenating slices before and after the element
			tasks = append(tasks[:i], tasks[i+1:]...)

			// Return 204 No Content (successful deletion)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// Task with given ID not found
	http.Error(w, "Task not found", http.StatusNotFound)
}
