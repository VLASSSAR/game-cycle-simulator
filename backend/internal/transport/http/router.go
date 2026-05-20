package http

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	handler := NewHandler()

	mux.HandleFunc("/api/health", handler.Health)
	mux.HandleFunc("/api/simulate", handler.Simulate)
	mux.HandleFunc("/api/optimize", handler.Optimize)
	mux.HandleFunc("/api/optimize/compare", handler.CompareOptimizationMethods)

	return enableCORS(mux)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
