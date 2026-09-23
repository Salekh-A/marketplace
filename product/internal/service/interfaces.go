package service

import (
	"context"

	"product/internal/model"
)

type ProductRepository interface {
	Create(
		ctx context.Context,
		product model.Product,
	) error

	GetByID(
		ctx context.Context,
		id string,
	) (model.Product, error)

	GetAll(
		ctx context.Context,
	) ([]model.Product, error)

	Update(
		ctx context.Context,
		product model.Product,
	) error

	DecreaseStock(
		ctx context.Context,
		id string,
		quantity int,
	) error

	Delete(
		ctx context.Context,
		id string,
	) error
}
