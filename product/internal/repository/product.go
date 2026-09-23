package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"product/internal/model"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(
	ctx context.Context,
	product model.Product,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO products (
			id,
			name,
			price,
			stock,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *ProductRepository) GetByID(
	ctx context.Context,
	id string,
) (model.Product, error) {
	var product model.Product

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			name,
			price,
			stock,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
		`,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	return product, err
}

func (r *ProductRepository) GetAll(
	ctx context.Context,
) ([]model.Product, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			name,
			price,
			stock,
			created_at,
			updated_at
		FROM products
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product

		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Update(
	ctx context.Context,
	product model.Product,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		UPDATE products
		SET name = $1,
			price = $2,
			stock = $3,
			updated_at = NOW()
		WHERE id = $4
		`,
		product.Name,
		product.Price,
		product.Stock,
		product.ID,
	)

	return err
}

func (r *ProductRepository) DecreaseStock(
	ctx context.Context,
	id string,
	quantity int,
) error {
	result, err := r.db.Exec(
		ctx,
		`
		UPDATE products
		SET
			stock = stock - $1,
			updated_at = NOW()
		WHERE id = $2
		  AND stock >= $1
		`,
		quantity,
		id,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("not enough product stock")
	}

	return nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id string,
) error {
	_, err := r.db.Exec(
		ctx,
		`DELETE FROM products WHERE id = $1`,
		id,
	)

	return err
}
