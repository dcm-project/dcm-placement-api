package service

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/opa"
)

type OpaService struct {
	validator *opa.Validator
}

func NewOpaService(validator *opa.Validator) *OpaService {
	return &OpaService{validator: validator}
}

func (s *OpaService) ValidateVmPlacement(ctx context.Context, env string, network string) (bool, error) {
	return s.validator.EvalPolicy(ctx, "subnet", map[string]string{
		"env":     env,
		"network": network,
	})
}
