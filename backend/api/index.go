//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// Vercel serverless entry point. main.go (used by docker-compose and the
// GitHub artifact) is untouched; this wraps the same router/middleware for
// Vercel's Go runtime. vercel.json rewrites every path to /api/index so
// routing.NewRouter() still dispatches on the original request path.
//

package handler

import (
	"net/http"
	"sync"

	"cartepro/database"
	"cartepro/routing"
)

var (
	initOnce sync.Once
	mux      http.Handler
)

func Handler(w http.ResponseWriter, r *http.Request) {
	initOnce.Do(func() {
		database.InitDatabase()
		database.SeedTestClient()
		database.SeedTestPartners()
		mux = routing.CorsMiddleware(routing.NewRouter())
	})
	mux.ServeHTTP(w, r)
}
