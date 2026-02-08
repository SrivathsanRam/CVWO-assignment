package router

import (
	"context"
	"net/http"
	"os"
	"strings"
	"strconv"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/routes"
	"github.com/go-chi/chi/v5"
)

func Setup() chi.Router {
	r := chi.NewRouter()
	r.Use(corsMiddleware)
	r.Use(authMiddleware)
	setUpRoutes(r)
	return r
}

func setUpRoutes(r chi.Router) {
	r.Group(routes.GetRoutes())
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "http://localhost:5173"
		}

		origin := r.Header.Get("Origin")

		// Check if origin is in allowed list (comma-separated)
		allowed := false
		for _, allowedOrigin := range splitOrigins(allowedOrigins) {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userIDHeader := r.Header.Get("X-User-ID")
        if userIDHeader != "" {
            if userID, err := strconv.Atoi(userIDHeader); err == nil {
                ctx := context.WithValue(r.Context(), "user_id", userID)
                r = r.WithContext(ctx)
            }
        }
        next.ServeHTTP(w, r)
    })
}


func splitOrigins(origins string) []string {
	var result []string
	for _, o := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
