package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		"postgres://postgres:postgres@localhost:5435/analytics?sslmode=disable",
	)
	require.NoError(t, err)

	err = db.Ping(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Close()
	})

	_, err = db.Exec(
		ctx,
		"TRUNCATE TABLE analytics_orders",
	)
	require.NoError(t, err)

	return db
}

func TestAnalyticsRepository_RecordOrder(t *testing.T) {
	db := setupTestDB(t)

	repository := NewAnalyticsRepository(db)

	createdAt := time.Now()

	err := repository.RecordOrder(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		2,
		100,
		"created",
		createdAt,
	)

	require.NoError(t, err)

	var count int

	err = db.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM analytics_orders WHERE id = $1",
		"11111111-1111-1111-1111-111111111111",
	).Scan(&count)

	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestAnalyticsRepository_GetRevenue(t *testing.T) {
	db := setupTestDB(t)

	repository := NewAnalyticsRepository(db)

	err := repository.RecordOrder(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		2,
		100,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	err = repository.RecordOrder(
		context.Background(),
		"22222222-2222-2222-2222-222222222222",
		"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		1,
		50,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	result, err := repository.GetRevenue(context.Background())

	require.NoError(t, err)
	require.Equal(t, 250.0, result.Revenue)
}

func TestAnalyticsRepository_GetOrders(t *testing.T) {
	db := setupTestDB(t)

	repository := NewAnalyticsRepository(db)

	err := repository.RecordOrder(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		1,
		100,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	err = repository.RecordOrder(
		context.Background(),
		"22222222-2222-2222-2222-222222222222",
		"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		1,
		200,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	err = repository.RecordOrder(
		context.Background(),
		"33333333-3333-3333-3333-333333333333",
		"cccccccc-cccc-cccc-cccc-cccccccccccc",
		1,
		300,
		"cancelled",
		time.Now(),
	)
	require.NoError(t, err)

	result, err := repository.GetOrders(context.Background())

	require.NoError(t, err)
	require.Equal(t, 2, result.Orders)
}

func TestAnalyticsRepository_GetAverageCheck(t *testing.T) {
	db := setupTestDB(t)

	repository := NewAnalyticsRepository(db)

	err := repository.RecordOrder(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		1,
		100,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	err = repository.RecordOrder(
		context.Background(),
		"22222222-2222-2222-2222-222222222222",
		"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		1,
		200,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	result, err := repository.GetAverageCheck(context.Background())

	require.NoError(t, err)
	require.Equal(t, 150.0, result.Average)
}

func TestAnalyticsRepository_GetTopProducts(t *testing.T) {
	db := setupTestDB(t)

	repository := NewAnalyticsRepository(db)

	err := repository.RecordOrder(
		context.Background(),
		"11111111-1111-1111-1111-111111111111",
		"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		2,
		100,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	err = repository.RecordOrder(
		context.Background(),
		"22222222-2222-2222-2222-222222222222",
		"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		1,
		50,
		"created",
		time.Now(),
	)
	require.NoError(t, err)

	result, err := repository.GetTopProducts(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Equal(t, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", result[0].ProductID)
	require.Equal(t, 200.0, result[0].Revenue)
	require.Equal(t, 1, result[0].Orders)
}
