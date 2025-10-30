package service

import (
	"context"
	"fmt"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/catalog"
	"github.com/dcm-project/dcm-placement-api/internal/handlers/v1alpha1/mappers"
	"github.com/dcm-project/dcm-placement-api/internal/opa"
	"github.com/dcm-project/dcm-placement-api/internal/provider"
	"github.com/dcm-project/dcm-placement-api/internal/store"
	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PlacementService struct {
	store           store.Store
	opa             *opa.Validator
	providerService *provider.Service
}

func NewPlacementService(store store.Store, opa *opa.Validator,
	providerService *provider.Service) *PlacementService {
	return &PlacementService{store: store, opa: opa, providerService: providerService}
}

func (s *PlacementService) CreateApplication(ctx context.Context, request *server.CreateApplicationJSONRequestBody, appID string) (*server.ApplicationResponse, error) {
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

	serviceType := request.Service
	var applicationID uuid.UUID
	if appID != "" {
		applicationID, _ = uuid.Parse(appID)
	} else {
		applicationID = uuid.New()
	}

	// Store in database post validation
	zones := s.opa.GetRequiredZones(result)
	if len(zones) == 0 {
		return nil, fmt.Errorf("no zones found")
	}

	appModel := model.Application{
		ID:            applicationID,
		Name:          request.Name,
		Service:       string(request.Service),
		Zones:         zones,
		Tier:          tier,
		DeploymentIDs: []string{},
	}

	app, err := s.store.Application().Create(ctx, appModel)
	if err != nil {
		return nil, err
	}

	// Deploy to provider service
	var deploymentIDs []string
	for _, zone := range zones {
		logger.Info("Creating deployment in Zone: ", "Zone: ", zone)
		var deploymentID string

		if serviceType == "webserver" {
			vm := catalog.GetCatalogVm(request.Service)
			deploymentID, err = s.providerService.CreateVMDeployment(ctx, request.Name, zone, vm, app.ID.String())
			if err != nil {
				// Rollback: delete already created deployments
				for _, id := range deploymentIDs {
					_ = s.providerService.DeleteDeployment(ctx, id)
				}
				_ = s.store.Application().Delete(ctx, app.ID)
				return nil, fmt.Errorf("failed to create VM deployment in zone %s: %w", zone, err)
			}
		} else if serviceType == "container" {
			containerApp := catalog.GetContainerApp()
			deploymentID, err = s.providerService.CreateContainerDeployment(ctx, request.Name, zone, containerApp, app.ID.String())
			if err != nil {
				// Rollback: delete already created deployments
				for _, id := range deploymentIDs {
					_ = s.providerService.DeleteDeployment(ctx, id)
				}
				_ = s.store.Application().Delete(ctx, app.ID)
				return nil, fmt.Errorf("failed to create container deployment in zone %s: %w", zone, err)
			}
		}

		deploymentIDs = append(deploymentIDs, deploymentID)
	}

	// Update application with deployment IDs
	app.DeploymentIDs = deploymentIDs
	app, err = s.store.Application().Update(ctx, *app)
	if err != nil {
		// Rollback: delete deployments
		for _, id := range deploymentIDs {
			_ = s.providerService.DeleteDeployment(ctx, id)
		}
		return nil, fmt.Errorf("failed to update application with deployment IDs: %w", err)
	}

	appService := string(request.Service)
	return &server.ApplicationResponse{
		Name:    &request.Name,
		Service: &appService,
		Tier:    &tier,
		Id:      &app.ID,
	}, nil
}

func (s *PlacementService) DeleteApplication(ctx context.Context, id uuid.UUID) (*server.ApplicationResponse, error) {
	logger := zap.S().Named("placement_service:delete_app")
	app, err := s.store.Application().Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Delete deployments from provider service
	for _, deploymentID := range app.DeploymentIDs {
		logger.Info("Deleting deployment: ", "DeploymentID: ", deploymentID)
		err = s.providerService.DeleteDeployment(ctx, deploymentID)
		if err != nil {
			logger.Warnw("Failed to delete deployment", "deploymentID", deploymentID, "error", err)
			// Continue deleting other deployments even if one fails
		}
	}

	// Delete app from database
	err = s.store.Application().Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return mappers.ApplicationToAPI(*app), nil
}
