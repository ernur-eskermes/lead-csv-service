package storage

import "github.com/ernur-eskermes/lead-csv-service/pkg/database/postgresql"

type Storage struct {
	Product *Product
}

func New(db postgresql.Client) *Storage {
	return &Storage{
		Product: NewProduct(db),
	}
}
