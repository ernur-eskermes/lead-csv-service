package restHandler

import (
	"errors"

	"github.com/ernur-eskermes/lead-csv-service/internal/core"
	appError "github.com/ernur-eskermes/lead-csv-service/pkg/app_error"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type Product struct {
	service ProductService
}

func NewProduct(service ProductService) *Product {
	return &Product{
		service: service,
	}
}

func (h *Product) initProductRoutes(api fiber.Router) {
	products := api.Group("/products")
	{
		products.Get("", h.getAllProducts)
		products.Post("", h.createProduct)
		products.Put(":id", h.updateProduct)
		products.Delete(":id", h.deleteProduct)
	}
}

// @Summary Get Products CSV
// @Tags products
// @Description Get products csv
// @ModuleID getAllProducts
// @Success 200
// @Produce  application/csv
// @Router /products [get]
func (h *Product) getAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAll(c.Context())
	if err != nil {
		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	b, err := gocsv.MarshalBytes(products)
	if err != nil {
		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	c.Attachment("products.csv")

	return c.Status(fiber.StatusOK).Send(b)
}

// @Summary Create Product
// @Tags products
// @Description Create products
// @ModuleID createProduct
// @Accept  json
// @Produce  json
// @Param input body core.CreateProductInput true "create product"
// @Success 201 {object} core.Product
// @Failure 400 {object} appError.Errors
// @Router /products [post]
func (h *Product) createProduct(c *fiber.Ctx) error {
	var inp core.CreateProductInput
	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := inp.Validate(); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	product, err := h.service.Create(c.Context(), inp)
	if err != nil {
		if errors.Is(err, core.ErrProductNameDuplicate) {
			return c.Status(fiber.StatusBadRequest).JSON(appError.Errors{core.ErrProductNameDuplicate})
		}

		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// @Summary Update Product
// @Tags products
// @Description Update products
// @ModuleID updateProduct
// @Accept  json
// @Produce  json
// @Param id path string true "product id"
// @Param input body core.UpdateProductInput true "update product"
// @Success 200
// @Failure 400 {object} appError.Errors
// @Failure 404 {object} response
// @Router /products/{id} [put]
func (h *Product) updateProduct(c *fiber.Ctx) error {
	var inp core.UpdateProductInput
	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := inp.Validate(); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	productID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{"id type unknown"})
	}

	if _, err = h.service.GetByID(c.Context(), productID); err != nil {
		if errors.Is(err, core.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response{err.Error()})
		}

		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	if err = h.service.Update(c.Context(), productID, inp); err != nil {
		if errors.Is(err, core.ErrProductNameDuplicate) {
			return c.Status(fiber.StatusBadRequest).JSON(appError.Errors{core.ErrProductNameDuplicate})
		}

		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusOK).Send(nil)
}

// @Summary Delete Product
// @Tags products
// @Description Delete product
// @ModuleID deleteProduct
// @Accept  json
// @Produce  json
// @Param id path string true "product id"
// @Success 204
// @Failure 404 {object} response
// @Router /products/{id} [delete]
func (h *Product) deleteProduct(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{"id type unknown"})
	}

	if err = h.service.Delete(c.Context(), productID); err != nil {
		if errors.Is(err, core.ErrProductNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(response{err.Error()})
		}

		log.Error(err)

		return c.Status(fiber.StatusInternalServerError).Send(nil)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
