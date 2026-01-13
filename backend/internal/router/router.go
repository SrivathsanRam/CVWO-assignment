package router

import (
    "net/http"
    "os"

    "github.com/SrivathsanRam/CVWO_project/backend/internal/routes"
    "github.com/go-chi/chi/v5"
)

func Setup() chi.Router {
    r := chi.NewRouter()
    r.Use(corsMiddleware)
    setUpRoutes(r)
    return r
}

func setUpRoutes(r chi.Router) {
    r.Group(routes.GetRoutes())
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get allowed origin from environment or use defaults
        allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
        if allowedOrigin == "" {
            // Default to localhost for development
            allowedOrigin = "http://localhost:5173"
        }

        origin := r.Header.Get("Origin")
        
        // Allow the specific origin or localhost for development
        if origin == allowedOrigin || origin == "http://localhost:5173" {
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