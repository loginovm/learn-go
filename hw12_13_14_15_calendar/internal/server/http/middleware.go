package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/logger"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := wrapResponseWriter(w)
		start := time.Now()
		next.ServeHTTP(wrapped, r)
		duration := time.Since(start).Milliseconds()

		ip := getUserIP(r)
		ua := getUserAgent(r)
		logger.Info(r.Context(),
			fmt.Sprintf("%s %s %s %s %d %d %s",
				ip,
				r.Method,
				r.URL.EscapedPath(),
				r.Proto,
				wrapped.status,
				duration,
				ua))
	})
}

func getUserIP(r *http.Request) string {
	var IPAddress string
	header := r.Header.Get("X-Forwarded-For")
	if header != "" {
		ips := strings.Split(header, ", ")
		if len(ips) > 0 {
			IPAddress = ips[len(ips)-1]
		}
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
		colonIdx := strings.IndexByte(IPAddress, ':')
		if colonIdx > 0 {
			IPAddress = IPAddress[:colonIdx]
		}
	}

	return IPAddress
}

func getUserAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}
