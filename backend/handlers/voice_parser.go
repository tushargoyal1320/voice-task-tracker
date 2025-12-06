/**
 * Voice Parser Package
 *
 * Contains AI-powered voice transcript parsing functions
 * Extracts task details from natural language input
 */
package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/**
 * getCurrentTimestamp returns current time in ISO format
 */
func getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

/**
 * ParseVoiceInput intelligently parses voice transcript to extract task details
 *
 * HTTP Method: POST
 * Endpoint: /parse-voice
 * Request Body: JSON with 'transcript' field
 * Response: JSON with parsed task fields (title, priority, dueDate, status)
 *
 * This function uses natural language processing to extract:
 * - Title: Main task description
 * - Priority: Keywords like "urgent", "high priority", "low priority"
 * - Due Date: Relative ("tomorrow", "next Monday") and absolute dates
 * - Status: Defaults to "To Do"
 */
func ParseVoiceInput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode the voice transcript from request
	var req VoiceParseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transcript := strings.TrimSpace(req.Transcript)
	if transcript == "" {
		http.Error(w, "Transcript cannot be empty", http.StatusBadRequest)
		return
	}

	// Parse the transcript
	response := VoiceParseResponse{
		Transcript: transcript,
		Status:     "To Do", // Default status
	}

	// Extract priority
	response.Priority = extractPriority(transcript)

	// Extract due date
	response.DueDate = extractDueDate(transcript)

	// Extract title (remove priority and date keywords)
	response.Title = extractTitle(transcript)

	// Send parsed response
	json.NewEncoder(w).Encode(response)
}

/**
 * extractPriority identifies priority level from transcript
 * Looks for keywords: urgent, critical, high, low priority, etc.
 */
func extractPriority(text string) string {
	lower := strings.ToLower(text)

	if strings.Contains(lower, "urgent") || strings.Contains(lower, "critical") {
		return "Urgent"
	}
	if strings.Contains(lower, "high priority") || strings.Contains(lower, "high") {
		return "High"
	}
	if strings.Contains(lower, "low priority") || strings.Contains(lower, "low") {
		return "Low"
	}

	return "Medium" // Default
}

/**
 * extractDueDate parses relative and absolute dates from transcript
 * Handles: tomorrow, today, next Monday, in 3 days, by Friday, etc.
 */
func extractDueDate(text string) string {
	lower := strings.ToLower(text)
	now := time.Now()

	// Tomorrow
	if strings.Contains(lower, "tomorrow") {
		return now.AddDate(0, 0, 1).Format("2006-01-02")
	}

	// Today
	if strings.Contains(lower, "today") {
		return now.Format("2006-01-02")
	}

	// Next Monday, Tuesday, etc.
	weekdays := map[string]time.Weekday{
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		"sunday":    time.Sunday,
	}

	for day, weekday := range weekdays {
		if strings.Contains(lower, "next "+day) {
			daysUntil := (int(weekday) - int(now.Weekday()) + 7) % 7
			if daysUntil == 0 {
				daysUntil = 7
			}
			return now.AddDate(0, 0, daysUntil).Format("2006-01-02")
		}
	}

	// "in X days"
	re := regexp.MustCompile(`in (\d+) days?`)
	matches := re.FindStringSubmatch(lower)
	if len(matches) > 1 {
		days, _ := strconv.Atoi(matches[1])
		return now.AddDate(0, 0, days).Format("2006-01-02")
	}

	return "" // No date found
}

/**
 * extractTitle removes date and priority keywords to get clean task title
 */
func extractTitle(text string) string {
	// Remove common command phrases
	title := text
	prefixes := []string{
		"create a task to ",
		"create task to ",
		"remind me to ",
		"add task to ",
		"add a task to ",
		"create a ",
		"add a ",
	}

	lower := strings.ToLower(title)
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			title = title[len(prefix):]
			break
		}
	}

	// Remove date-related phrases
	datePatterns := []string{
		` by tomorrow.*`,
		` tomorrow.*`,
		` today.*`,
		` by next \w+.*`,
		` next \w+.*`,
		` in \d+ days.*`,
		` by \d+.*`,
	}

	for _, pattern := range datePatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		title = re.ReplaceAllString(title, "")
	}

	// Remove priority phrases
	priorityPatterns := []string{
		` it'?s urgent`,
		` it'?s high priority`,
		` it'?s low priority`,
		` high priority`,
		` low priority`,
		` urgent`,
		` critical`,
	}

	for _, pattern := range priorityPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		title = re.ReplaceAllString(title, "")
	}

	return strings.TrimSpace(title)
}

/**
 * SearchTasks filters tasks based on query parameters
 *
 * HTTP Method: GET
 * Endpoint: /tasks/search?q=query&status=status&priority=priority
 * Query Parameters:
 *   - q: Search in title and description
 *   - status: Filter by status
 *   - priority: Filter by priority
 * Response: JSON array of matching tasks
 */
func SearchTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	query := strings.ToLower(r.URL.Query().Get("q"))
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	mu.Lock()
	defer mu.Unlock()

	// Filter tasks
	var filtered []Task
	for _, task := range tasks {
		// Text search
		if query != "" {
			titleMatch := strings.Contains(strings.ToLower(task.Title), query)
			descMatch := strings.Contains(strings.ToLower(task.Description), query)
			if !titleMatch && !descMatch {
				continue
			}
		}

		// Status filter
		if status != "" && task.Status != status {
			continue
		}

		// Priority filter
		if priority != "" && task.Priority != priority {
			continue
		}

		filtered = append(filtered, task)
	}

	json.NewEncoder(w).Encode(filtered)
}
