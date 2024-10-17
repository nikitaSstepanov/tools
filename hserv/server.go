package hserv

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikitaSstepanov/tools/sl"
)

type Config struct {
	Url             string        `yaml:"url"             env:"SERVER_URL"              env-default:":80"`
	ReadTimeout     time.Duration `yaml:"readTimeout"     env:"SERVER_READ_TIMEOUT"     env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"writeTimeout"    env:"SERVER_WRITE_TIMEOUT"    env-default:"5s"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout" env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"5s"`
}

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, cfg *Config) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Addr:         cfg.Url,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: cfg.ShutdownTimeout,
	}

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown(ctx context.Context) error {
	log := sl.L(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("signal: " + s.String())
	case err := <-s.Notify():
		log.Error("httpServer.Notify:" + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
