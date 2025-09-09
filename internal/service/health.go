package service

import (
	"context"
)

type HealthService struct {
}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Health(ctx context.Context) error {
	return nil
}
