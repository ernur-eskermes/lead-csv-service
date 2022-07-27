package core

import (
	"errors"

	appError "github.com/ernur-eskermes/lead-csv-service/pkg/app_error"
)

var (
	ErrProductNameDuplicate = appError.New("name", "product name already exists")

	ErrProductNotFound = errors.New("product not found")
)

type Product struct {
	ID    int    `json:"id" csv:"-"`
	Name  string `json:"name" csv:"PRODUCT NAME"`
	Price int    `json:"price" csv:"PRICE"`
}

type CreateProductInput struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"number,min=0"`
}

func (c *CreateProductInput) Validate() appError.Errors {
	return validateStruct(c)
}

type UpdateProductInput struct {
	Name  string `json:"name" validate:"required"`
	Price int    `json:"price" validate:"number,min=0"`
}

func (c *UpdateProductInput) Validate() appError.Errors {
	return validateStruct(c)
}
