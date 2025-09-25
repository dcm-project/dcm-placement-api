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

func (s *PlacementService) CreateApplication(ctx context.Context, request *server.CreateApplicationJSONRequestBody) (*server.Application, error) {
	logger := zap.S().Named("placement_service:create_app")

	// OPA validation:
	tier := 2
	if request.Tier != nil {
		tier = *request.Tier
	}
	logger.Info("Evaluating policy: ", "Tier: ", fmt.Sprintf("%d", tier))
	result, err := s.opa.EvalPolicy(ctx, fmt.Sprintf("tier%d", tier), map[string]string{
		"name": request.Name,
	})
	if err != nil {
		return nil, err
	}

	logger.Info("OPA validation result: ", "Result: ", result)

	if !s.opa.IsInputValid(result) {
		//logger.Warn("Invalid input: ", "Input: ", request)
	}

	if !s.opa.IsOutputValid(result) {
		return nil, fmt.Errorf("invalid output")
	}

	zones := s.opa.GetOutputZones(result)

	// Store in database post validation
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
