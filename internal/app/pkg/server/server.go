package server

import (
	"context"
	"fmt"
	"net/http"
	"otp_api/conf"
	"otp_api/internal/app/controller"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func NewServer(cfg *conf.Config, handler *controller.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Sever.Host, cfg.Sever.Port),
		Handler:      handler.Setup(),
		ReadTimeout:  time.Duration(cfg.Sever.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Sever.WriteTimeout) * time.Second,
	}}

}
