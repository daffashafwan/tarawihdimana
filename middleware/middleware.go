package middleware

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	// Create a rate limiter that allows 10 requests per second with burst capability
	rateLmInt, err := strconv.ParseInt(os.Getenv("RATE_LIMIT_MAX"), 10, 64)
	if err != nil {
		rateLmInt = 10
	}
	limiter := ratelimit.NewBucket(time.Second, rateLmInt)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.TakeAvailable(1) == 0 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CorsMiddleware(next http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins(GetAllowedOrigins())
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return handlers.CORS(headersOk, originsOk, methodsOk)(next)
}

func GetAllowedOrigins() []string{
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	log.Printf("Allowed Origins: %v\n", allowedOrigins)
	if allowedOrigins == "" {
		return []string{"*"}
	}

	return strings.Split(allowedOrigins, ",")
}