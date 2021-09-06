package restserver

import (
	"time"

	"github.com/brandenc40/safer"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app   *fiber.App
	safer *safer.Client
}

func NewServer() *Server {
	s := Server{
		app: fiber.New(fiber.Config{
			ReadTimeout: 5 * time.Second,
			AppName:     "SAFER Scraper",
		}),
		safer: safer.NewClient(),
	}
	s.app.Get("/snapshot/dot/:value", func(ctx *fiber.Ctx) error {
		snapshot, err := s.safer.GetCompanyByDOTNumber(ctx.Params("value"))
		if err != nil {
			return err
		}
		return ctx.JSON(snapshot)
	})
	s.app.Get("/snapshot/mcmx/:value", func(ctx *fiber.Ctx) error {
		snapshot, err := s.safer.GetCompanyByMCMX(ctx.Params("value"))
		if err != nil {
			return err
		}
		return ctx.JSON(snapshot)
	})
	s.app.Get("/search/:value", func(ctx *fiber.Ctx) error {
		snapshot, err := s.safer.SearchCompaniesByName(ctx.Params("value"))
		if err != nil {
			return err
		}
		return ctx.JSON(snapshot)
	})
	return &s
}

func (s *Server) StartServer(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) StopServer() error {
	return s.app.Shutdown()
}
