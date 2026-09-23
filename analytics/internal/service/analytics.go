package service

import (
	"context"
	"encoding/json"
	"time"

	"analytics/internal/model"
	"analytics/internal/repository"
)

type AnalyticsService struct {
	repository AnalyticsRepository
	redis      Cache
}

func NewAnalyticsService(
	repository AnalyticsRepository,
	redis Cache,
) *AnalyticsService {
	return &AnalyticsService{
		repository: repository,
		redis:      redis,
	}
}

func (s *AnalyticsService) GetRevenue(
	ctx context.Context,
) (model.Revenue, error) {
	var result model.Revenue

	if err := s.getCache(
		ctx,
		"analytics:revenue",
		&result,
	); err == nil {
		return result, nil
	}

	result, err := s.repository.GetRevenue(ctx)
	if err != nil {
		return result, err
	}

	s.setCache(
		ctx,
		"analytics:revenue",
		result,
	)

	return result, nil
}

func (s *AnalyticsService) GetOrders(
	ctx context.Context,
) (model.OrderStats, error) {
	var result model.OrderStats

	if err := s.getCache(
		ctx,
		"analytics:orders",
		&result,
	); err == nil {
		return result, nil
	}

	result, err := s.repository.GetOrders(ctx)
	if err != nil {
		return result, err
	}

	s.setCache(
		ctx,
		"analytics:orders",
		result,
	)

	return result, nil
}

func (s *AnalyticsService) GetAverageCheck(
	ctx context.Context,
) (model.AverageCheck, error) {
	var result model.AverageCheck

	if err := s.getCache(
		ctx,
		"analytics:average-check",
		&result,
	); err == nil {
		return result, nil
	}

	result, err := s.repository.GetAverageCheck(ctx)
	if err != nil {
		return result, err
	}

	s.setCache(
		ctx,
		"analytics:average-check",
		result,
	)

	return result, nil
}

func (s *AnalyticsService) GetTopProducts(
	ctx context.Context,
) ([]model.ProductStats, error) {
	var result []model.ProductStats

	if err := s.getCache(
		ctx,
		"analytics:top-products",
		&result,
	); err == nil {
		return result, nil
	}

	result, err := s.repository.GetTopProducts(ctx)
	if err != nil {
		return nil, err
	}

	s.setCache(
		ctx,
		"analytics:top-products",
		result,
	)

	return result, nil
}

func (s *AnalyticsService) getCache(
	ctx context.Context,
	key string,
	result interface{},
) error {
	data, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal(
		[]byte(data),
		result,
	)
}

func (s *AnalyticsService) setCache(
	ctx context.Context,
	key string,
	data interface{},
) {
	value, err := json.Marshal(data)
	if err != nil {
		return
	}

	_ = s.redis.Set(
		ctx,
		key,
		value,
		time.Minute,
	).Err()
}

func (s *AnalyticsService) RecordOrder(
	ctx context.Context,
	id string,
	productID string,
	quantity int,
	price float64,
	status string,
	createdAt time.Time,
) error {
	err := s.repository.RecordOrder(
		ctx,
		id,
		productID,
		quantity,
		price,
		status,
		createdAt,
	)
	if err != nil {
		return err
	}

	cacheKeys := []string{
		"analytics:revenue",
		"analytics:orders",
		"analytics:average-check",
		"analytics:top-products",
	}

	for _, key := range cacheKeys {
		if err := s.redis.Del(ctx, key).Err(); err != nil {
			return err
		}
	}

	return nil
}

var _ AnalyticsRepository = (*repository.AnalyticsRepository)(nil)
