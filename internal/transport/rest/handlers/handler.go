package restHandler

import (
	"context"
	"time"

	"github.com/ernur-eskermes/lead-csv-service/internal/core"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]core.Product, error)
	Create(ctx context.Context, inp core.CreateProductInput) (core.Product, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, inp core.UpdateProductInput) error
	GetByID(ctx context.Context, id int) (core.Product, error)
}

type Deps struct {
	ProductService ProductService
}

type Handler struct {
	productHandler *Product
}

func New(deps Deps) *Handler {
	return &Handler{
		productHandler: NewProduct(deps.ProductService),
	}
}

func (h *Handler) InitRouter(app *fiber.App) {
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Almaty",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)
	h.initAPI(app)
}

func (h *Handler) initAPI(app fiber.Router) {
	api := app.Group("/api")
	{
		h.productHandler.initProductRoutes(api)
	}
}

type response struct {
	Message string `json:"message"`
}
