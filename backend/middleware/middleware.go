/**
 * Middleware Package
 *
 * Contains HTTP middleware functions that intercept and process requests
 * Middleware runs before the actual handler functions execute
 */
package middleware

import (
	"log"
	"net/http"
)

/**
 * LoggingMiddleware logs details of each incoming HTTP request
 *
 * Purpose: Debugging and monitoring - helps track API usage and diagnose issues
 * Logs: HTTP method (GET/POST/PUT/DELETE) and URL path
 *
 * Example output: "Received request: POST /tasks"
 */
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request method and URL path
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		// Pass control to the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

/**
 * CORSMiddleware adds Cross-Origin Resource Sharing (CORS) headers
 *
 * Purpose: Allows frontend (localhost:3000) to make requests to backend (localhost:8080)
 * Without CORS, browsers block requests between different origins for security
 *
 * CORS Headers:
 * - Access-Control-Allow-Origin: Specifies allowed origins (* = all origins)
 * - Access-Control-Allow-Methods: Lists allowed HTTP methods
 * - Access-Control-Allow-Headers: Lists allowed request headers
 *
 * OPTIONS Requests (Preflight):
 * Browsers send OPTIONS request before actual request to check permissions
 * This middleware responds to OPTIONS with 200 OK to approve the request
 */
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Only handle OPTIONS preflight here
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

/**
 * AuthMiddleware checks for authentication before allowing access
 *
 * Purpose: Protect routes that require authentication
 * Currently checks for Authorization header presence (placeholder implementation)
 *
 * Note: This is a basic example. In production, implement proper token validation:
 * - JWT token verification
 * - Session validation
 * - API key authentication
 *
 * Usage: Apply to specific routes that need protection
 * Example: router.Handle("/protected-route", AuthMiddleware(handler))
 */
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for Authorization header in request
		token := r.Header.Get("Authorization")
		if token == "" {
			// No token provided - deny access
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO: Add actual token validation logic here
		// - Verify JWT signature
		// - Check token expiration
		// - Validate user permissions

		// Pass control to the next handler if authenticated
		next.ServeHTTP(w, r)
	})
}
