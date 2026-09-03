//
// EPITECH PROJECT, 2026
// CartePro
// File description:
// cors
//

package routing

import (
	"net/http"
	"os"
	"slices"
	"strings"
)

//CorsMiddleware allows the frontend (a different origin) to call this API with credentials.
//It echoes back Access-Control-Allow-Origin only for an allow-listed origin (required
//alongside Allow-Credentials, since browsers reject "*" when credentials are used), and answers
//OPTIONS preflight requests directly since they aren't registered on any route in router.go.
//FRONTEND_ORIGIN may be a single origin or a comma-separated list (e.g. a prod domain plus
//Vercel preview URLs).
func CorsMiddleware(next http.Handler) http.Handler {
	var allowedOrigins []string
	for _, o := range strings.Split(os.Getenv("FRONTEND_ORIGIN"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowedOrigins = append(allowedOrigins, o)
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && slices.Contains(allowedOrigins, origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin) //echo back the origin only if it matches the allowed origin
			w.Header().Set("Access-Control-Allow-Credentials", "true") //allow credentials (cookies, authorization headers, etc.)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") //allow all methods, including OPTIONS for preflight requests
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type") //allow Content-Type header for JSON requests
			w.Header().Set("Vary", "Origin") //indicate that the response varies based on the Origin header, so caches don't serve the wrong response to other origins
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
