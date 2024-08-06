package internalhttp

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	addr   string
	logger *logger.Logger
	srv    *http.Server
}

type Application interface { // TODO
}

func NewServer(addr string, logger *logger.Logger, _ Application) *Server {
	return &Server{
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "hello-world")
		w.WriteHeader(http.StatusOK)
	})
	s.AddLogger(r)
	r.Use(loggingMiddleware)

	server := &http.Server{
		Addr:              s.addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	s.srv = server
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping calendar...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return s.srv.Shutdown(ctx)
}

// AddLogger adds logger to context
// so that it is available in all http handlers.
func (s *Server) AddLogger(r *mux.Router) {
	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithLogger(r.Context(), s.logger)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
	r.Use(handler)
}
