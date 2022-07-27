package rest

import (
	"fmt"

	restHandler "github.com/ernur-eskermes/lead-csv-service/internal/transport/rest/handlers"

	"github.com/ernur-eskermes/lead-csv-service/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	httpServer *fiber.App
	Handlers   *restHandler.Handler
}

func NewServer(cfg *config.Config, handlers *restHandler.Handler) *Server {
	app := fiber.New(fiber.Config{
		WriteTimeout: cfg.HTTP.WriteTimeout,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		BodyLimit:    cfg.HTTP.MaxHeaderMegabytes << 20,
	})
	handlers.InitRouter(app)

	return &Server{httpServer: app, Handlers: handlers}
}

func (s *Server) ListenAndServe(port int) error {
	return s.httpServer.Listen(fmt.Sprintf(":%d", port))
}

func (s *Server) Stop() error {
	return s.httpServer.Shutdown()
}
