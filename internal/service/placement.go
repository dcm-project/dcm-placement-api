package service

import (
	"context"
	"fmt"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/catalog"
	"github.com/dcm-project/dcm-placement-api/internal/deploy"
	"github.com/dcm-project/dcm-placement-api/internal/handlers/v1alpha1/mappers"
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

func (s *PlacementService) CreateApplication(ctx context.Context, request *server.CreateApplicationJSONRequestBody) (*server.Application, error) {
	logger := zap.S().Named("placement_service:create_app")

	// OPA validation:
	tier := 2
	if request.Tier != nil {
		tier = *request.Tier
	}
	logger.Info("Evaluating policy: ", "Tier: ", fmt.Sprintf("%d", tier))
	result, err := s.opa.EvalTierPolicy(ctx, tier, request.Name, request.Zones)
	if err != nil {
		return nil, err
	}

	logger.Info("OPA validation result: ", "Result: ", result)

	if !s.opa.IsValid(result) {
		failures := s.opa.GetFailures(result)
		if len(failures) > 0 {
			return nil, fmt.Errorf("validation failed: %v", failures)
		}
		return nil, fmt.Errorf("input validation failed")
	}

	// Store in database post validation
	zones := s.opa.GetRequiredZones(result)
	app, err := s.store.Application().Create(ctx, model.Application{
		ID:      uuid.New(),
		Name:    request.Name,
		Service: string(request.Service),
		Zones:   zones,
		Tier:    tier,
	})
	if err != nil {
		return nil, err
	}

	// Deploy VMs
	for _, zone := range zones {
		logger.Info("Created service in Zone: ", "Zone: ", zone)
		vm := catalog.GetCatalogVm(request.Service)
		err = s.deploy.DeployVM(ctx, fmt.Sprintf("%s-%s", request.Name, app.ID.String()), vm, zone)
		if err != nil {
			return nil, err
		}
	}

	return &server.Application{
		Name:    request.Name,
		Service: request.Service,
		Tier:    &tier,
		Id:      &app.ID,
	}, nil
}

func (s *PlacementService) DeleteApplication(ctx context.Context, id uuid.UUID) (*server.Application, error) {
	logger := zap.S().Named("placement_service:delete_app")
	app, err := s.store.Application().Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete VMs from zones
	for _, zone := range app.Zones {
		logger.Info("Removing service from Zone: ", "Zone: ", zone)
		err = s.deploy.DeleteVM(ctx, fmt.Sprintf("%s-%s", app.Name, id), zone)
		if err != nil {
			return nil, err
		}
	}

	// Delete app from database
	err = s.store.Application().Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.ApplicationToAPI(*app), nil
}
