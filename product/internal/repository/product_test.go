package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"product/internal/model"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5434/products?sslmode=disable",
	)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		t.Fatalf("failed to ping database: %v", err)
	}

	return db
}

func createTestProduct(t *testing.T, repository *ProductRepository) model.Product {
	t.Helper()

	now := time.Now()

	product := model.Product{
		ID:        uuid.New().String(),
		Name:      "Test Product",
		Price:     1000,
		Stock:     10,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repository.Create(context.Background(), product); err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return product
}

func TestProductRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)
	ctx := context.Background()

	product := model.Product{
		ID:        uuid.New().String(),
		Name:      "iPhone",
		Price:     999.99,
		Stock:     10,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			ctx,
			"DELETE FROM products WHERE id = $1",
			product.ID,
		)
	})

	err := repository.Create(ctx, product)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var count int

	err = db.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM products WHERE id = $1",
		product.ID,
	).Scan(&count)
	if err != nil {
		t.Fatalf("failed to verify product: %v", err)
	}

	if count != 1 {
		t.Fatalf(
			"expected product count %d, got %d",
			1,
			count,
		)
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	product := createTestProduct(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM products WHERE id = $1",
			product.ID,
		)
	})

	result, err := repository.GetByID(
		context.Background(),
		product.ID,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != product.ID {
		t.Fatalf(
			"expected ID %q, got %q",
			product.ID,
			result.ID,
		)
	}

	if result.Name != product.Name {
		t.Fatalf(
			"expected name %q, got %q",
			product.Name,
			result.Name,
		)
	}

	if result.Price != product.Price {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			product.Price,
			result.Price,
		)
	}

	if result.Stock != product.Stock {
		t.Fatalf(
			"expected stock %d, got %d",
			product.Stock,
			result.Stock,
		)
	}
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	_, err := repository.GetByID(
		context.Background(),
		uuid.New().String(),
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestProductRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	product1 := createTestProduct(t, repository)
	product2 := createTestProduct(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM products WHERE id IN ($1, $2)",
			product1.ID,
			product2.ID,
		)
	})

	products, err := repository.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found1 := false
	found2 := false

	for _, product := range products {
		if product.ID == product1.ID {
			found1 = true
		}

		if product.ID == product2.ID {
			found2 = true
		}
	}

	if !found1 {
		t.Fatalf("product %q not found", product1.ID)
	}

	if !found2 {
		t.Fatalf("product %q not found", product2.ID)
	}
}

func TestProductRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	product := createTestProduct(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM products WHERE id = $1",
			product.ID,
		)
	})

	product.Name = "Updated Product"
	product.Price = 1500
	product.Stock = 20

	err := repository.Update(
		context.Background(),
		product,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repository.GetByID(
		context.Background(),
		product.ID,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "Updated Product" {
		t.Fatalf(
			"expected name %q, got %q",
			"Updated Product",
			result.Name,
		)
	}

	if result.Price != 1500 {
		t.Fatalf(
			"expected price %.2f, got %.2f",
			1500.0,
			result.Price,
		)
	}

	if result.Stock != 20 {
		t.Fatalf(
			"expected stock %d, got %d",
			20,
			result.Stock,
		)
	}
}

func TestProductRepository_DecreaseStock(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	product := createTestProduct(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM products WHERE id = $1",
			product.ID,
		)
	})

	err := repository.DecreaseStock(
		context.Background(),
		product.ID,
		3,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := repository.GetByID(
		context.Background(),
		product.ID,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Stock != 7 {
		t.Fatalf(
			"expected stock %d, got %d",
			7,
			result.Stock,
		)
	}
}

func TestProductRepository_DecreaseStock_NotEnoughStock(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	product := createTestProduct(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM products WHERE id = $1",
			product.ID,
		)
	})

	err := repository.DecreaseStock(
		context.Background(),
		product.ID,
		100,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "not enough product stock" {
		t.Fatalf(
			"expected error %q, got %q",
			"not enough product stock",
			err.Error(),
		)
	}
}

func TestProductRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewProductRepository(db)

	product := createTestProduct(t, repository)

	err := repository.Delete(
		context.Background(),
		product.ID,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repository.GetByID(
		context.Background(),
		product.ID,
	)

	if err == nil {
		t.Fatal("expected product to be deleted")
	}
}
