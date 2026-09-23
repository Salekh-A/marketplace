package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"order/internal/model"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5433/orders?sslmode=disable",
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

func createTestOrder(
	t *testing.T,
	repository *OrderRepository,
) model.Order {
	t.Helper()

	now := time.Now()

	order := model.Order{
		ID:        uuid.New().String(),
		ProductID: uuid.New().String(),
		Quantity:  2,
		Price:     999.99,
		Status:    "created",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := repository.Create(
		context.Background(),
		order,
	); err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}

	return order
}

func TestOrderRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewOrderRepository(db)

	order := createTestOrder(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM orders WHERE id = $1",
			order.ID,
		)
	})

	result, err := repository.GetByID(
		context.Background(),
		order.ID,
	)
	if err != nil {
		t.Fatalf("failed to get created order: %v", err)
	}

	if result.ID != order.ID {
		t.Errorf(
			"expected id %q, got %q",
			order.ID,
			result.ID,
		)
	}

	if result.ProductID != order.ProductID {
		t.Errorf(
			"expected product_id %q, got %q",
			order.ProductID,
			result.ProductID,
		)
	}

	if result.Quantity != order.Quantity {
		t.Errorf(
			"expected quantity %d, got %d",
			order.Quantity,
			result.Quantity,
		)
	}

	if result.Price != order.Price {
		t.Errorf(
			"expected price %v, got %v",
			order.Price,
			result.Price,
		)
	}

	if result.Status != order.Status {
		t.Errorf(
			"expected status %q, got %q",
			order.Status,
			result.Status,
		)
	}
}

func TestOrderRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewOrderRepository(db)

	order := createTestOrder(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM orders WHERE id = $1",
			order.ID,
		)
	})

	result, err := repository.GetByID(
		context.Background(),
		order.ID,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != order.ID {
		t.Errorf(
			"expected id %q, got %q",
			order.ID,
			result.ID,
		)
	}

	if result.ProductID != order.ProductID {
		t.Errorf(
			"expected product_id %q, got %q",
			order.ProductID,
			result.ProductID,
		)
	}

	if result.Quantity != order.Quantity {
		t.Errorf(
			"expected quantity %d, got %d",
			order.Quantity,
			result.Quantity,
		)
	}

	if result.Price != order.Price {
		t.Errorf(
			"expected price %v, got %v",
			order.Price,
			result.Price,
		)
	}

	if result.Status != order.Status {
		t.Errorf(
			"expected status %q, got %q",
			order.Status,
			result.Status,
		)
	}
}

func TestOrderRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewOrderRepository(db)

	_, err := repository.GetByID(
		context.Background(),
		uuid.New().String(),
	)

	if err == nil {
		t.Fatal("expected error for non-existent order")
	}
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewOrderRepository(db)

	order := createTestOrder(t, repository)

	t.Cleanup(func() {
		_, _ = db.Exec(
			context.Background(),
			"DELETE FROM orders WHERE id = $1",
			order.ID,
		)
	})

	err := repository.UpdateStatus(
		context.Background(),
		order.ID,
		"cancelled",
	)
	if err != nil {
		t.Fatalf("failed to update status: %v", err)
	}

	result, err := repository.GetByID(
		context.Background(),
		order.ID,
	)
	if err != nil {
		t.Fatalf("failed to get updated order: %v", err)
	}

	if result.Status != "cancelled" {
		t.Errorf(
			"expected status cancelled, got %q",
			result.Status,
		)
	}
}

func TestOrderRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repository := NewOrderRepository(db)

	order := createTestOrder(t, repository)

	err := repository.Delete(
		context.Background(),
		order.ID,
	)
	if err != nil {
		t.Fatalf("failed to delete order: %v", err)
	}

	_, err = repository.GetByID(
		context.Background(),
		order.ID,
	)
	if err == nil {
		t.Fatal("expected error after deleting order")
	}
}
