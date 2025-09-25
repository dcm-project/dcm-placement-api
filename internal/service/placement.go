package service

import (
	"context"
	"fmt"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/catalog"
	"github.com/dcm-project/dcm-placement-api/internal/deploy"
	"github.com/dcm-project/dcm-placement-api/internal/opa"
	"github.com/dcm-project/dcm-placement-api/internal/store"
	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PlacementService struct {
	store  store.Store
	opa    *opa.Validator
	deploy *deploy.DeployService
}

func NewPlacementService(store store.Store, opa *opa.Validator, deploy *deploy.DeployService) *PlacementService {
	return &PlacementService{store: store, opa: opa, deploy: deploy}
}

func (s *PlacementService) CreateApplication(ctx context.Context, request *server.CreateApplicationJSONRequestBody) error {
	logger := zap.S().Named("placement_service:create_app")

	// Handle tier field with default value
	tier := "tier1"
	if request.Tier != nil {
		tier = *request.Tier
	}

	// OPA validation:
	result, err := s.opa.EvalPolicy(ctx, tier, map[string]string{
		"name": request.Name,
		"tier": "1",
	})
	if err != nil {
		return err
	}

	logger.Info("OPA validation result", "Result", result)

	if !s.opa.IsInputValid(result) {
		logger.Warn("Invalid input", "Input", request)
	}

	if !s.opa.IsOutputValid(result) {
		return fmt.Errorf("invalid output")
	}

	zones := s.opa.GetOutputZones(result)
	for _, zone := range zones {
		logger.Info("Created Application in Zone", "Zone", zone)
		vm := catalog.GetCatalogVm(request.Service)
		err = s.deploy.DeployVM(ctx, request.Name, vm, zone)
		if err != nil {
			return err
		}
	}

	// Store in database post validation
	_, err = s.store.Application().Create(ctx, model.Application{
		ID:      uuid.New(),
		Name:    request.Name,
		Service: string(request.Service),
		Zones:   zones,
		Tier:    tier,
	})
	if err != nil {
		return err
	}

	return nil
}
