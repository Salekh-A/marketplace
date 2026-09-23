package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"product/internal/model"
)

type ProductService struct {
	repository ProductRepository
}

func NewProductService(
	repository ProductRepository,
) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) Create(
	ctx context.Context,
	name string,
	price float64,
	stock int,
) (model.Product, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return model.Product{}, errors.New("name is required")
	}

	if price < 0 {
		return model.Product{}, errors.New("price cannot be negative")
	}

	if stock < 0 {
		return model.Product{}, errors.New("stock cannot be negative")
	}

	now := time.Now()

	product := model.Product{
		ID:        uuid.New().String(),
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(ctx, product); err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *ProductService) GetByID(
	ctx context.Context,
	id string,
) (model.Product, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *ProductService) GetAll(
	ctx context.Context,
) ([]model.Product, error) {
	return s.repository.GetAll(ctx)
}

func (s *ProductService) Update(
	ctx context.Context,
	id string,
	name string,
	price float64,
	stock int,
) (model.Product, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return model.Product{}, errors.New("name is required")
	}

	if price < 0 {
		return model.Product{}, errors.New("price cannot be negative")
	}

	if stock < 0 {
		return model.Product{}, errors.New("stock cannot be negative")
	}

	product := model.Product{
		ID:    id,
		Name:  name,
		Price: price,
		Stock: stock,
	}

	if err := s.repository.Update(ctx, product); err != nil {
		return model.Product{}, err
	}

	return s.repository.GetByID(ctx, id)
}

func (s *ProductService) DecreaseStock(
	ctx context.Context,
	id string,
	quantity int,
) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	return s.repository.DecreaseStock(
		ctx,
		id,
		quantity,
	)
}

func (s *ProductService) Delete(
	ctx context.Context,
	id string,
) error {
	return s.repository.Delete(ctx, id)
}
