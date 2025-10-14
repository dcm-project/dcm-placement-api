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
	store            store.Store
	opa              *opa.Validator
	deploy           *deploy.DeployService
	containerService *deploy.ContainerService
}

func NewPlacementService(store store.Store, opa *opa.Validator,
	deploy *deploy.DeployService, containerService *deploy.ContainerService) *PlacementService {
	return &PlacementService{store: store, opa: opa, deploy: deploy, containerService: containerService}
}

func (s *PlacementService) CreateApplication(ctx context.Context, request *server.CreateApplicationJSONRequestBody, appID string) (*server.ApplicationResponse, error) {
	logger := zap.S().Named("placement_service:create_app")

	// OPA validation:
	tier := 2
	if request.Tier != nil {
		tier = *request.Tier
	}
	logger.Info("Policy Validation")
	zones, err := s.opaValidation(ctx, tier, request.Name, request.Zones)
	if err != nil {
		return nil, err
	}

	logger.Info("Post Policy Validation - Saving request in Database...")
	serviceType := request.Service
	var applicationID uuid.UUID
	if appID != "" {
		applicationID, _ = uuid.Parse(appID)
	} else {
		applicationID = uuid.New()
	}
	// Store in database post validation
	appModel := model.Application{
		ID:      applicationID,
		Name:    request.Name,
		Service: string(request.Service),
		Zones:   zones,
		Tier:    tier,
	}

	app, err := s.store.Application().Create(ctx, appModel)
	if err != nil {
		return nil, err
	}

	logger.Info("Deploying Application...")
	err = s.deployApplication(ctx, serviceType, request.Name, applicationID.String(), zones)
	if err != nil {
		return nil, err
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

	// Delete VMs from zones
	for _, zone := range app.Zones {
		logger.Info("Removing service from Zone: ", "Zone: ", zone)
		err = s.deploy.DeleteVM(ctx, id.String(), zone)
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

func (s *PlacementService) UpdateApplication(ctx context.Context, requestID string) (*server.ApplicationResponse, error) {
	logger := zap.S().Named("placement_service:update_app")
	logger.Info("Updating Application. ", "Application: ", requestID)

	// check id exist in database
	logger.Info("Retrieving Application from database...")
	application, err := s.store.Application().Get(ctx, uuid.MustParse(requestID))
	if err != nil {
		return &server.ApplicationResponse{}, err
	}

	// OPA Validation
	logger.Info("Policy Validation...")
	zones := []string(application.Zones)
	validatedZones, err := s.opaValidation(ctx, application.Tier, application.Name, &zones)
	if err != nil {
		return &server.ApplicationResponse{}, err
	}

	logger.Info("Deploying Application")
	err = s.deployApplication(ctx, server.ApplicationService(application.Service), application.Name, application.ID.String(), validatedZones)
	if err != nil {
		return &server.ApplicationResponse{}, err
	}

	updateApplication := server.ApplicationResponse{
		Name:    &application.Name,
		Service: &application.Service,
		Tier:    &application.Tier,
		Zones:   &validatedZones,
	}
	logger.Info("Successfully updated application. ", "Application: ", application.ID.String())
	return &updateApplication, nil
}

func (s *PlacementService) opaValidation(ctx context.Context, tier int, name string, zones *[]string) ([]string, error) {
	logger := zap.S().Named("placement_service:get_app")

	logger.Info("Evaluating policy: ", "Tier: ", fmt.Sprintf("%d", tier))
	result, err := s.opa.EvalTierPolicy(ctx, tier, name, zones)
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
	validatedZones := s.opa.GetRequiredZones(result)
	return validatedZones, nil
}

func (s *PlacementService) deployApplication(ctx context.Context, serviceType server.ApplicationService, appName, appID string, zones []string) error {
	logger := zap.S().Named("placement_service:deploy_app")

	for _, zone := range zones {
		logger.Info("Created service in Zone: ", "Zone: ", zone)
		// Deploy VMs
		if serviceType == "webserver" {
			vm := catalog.GetCatalogVm(serviceType)
			err := s.deploy.DeployVM(ctx, appName, vm, zone, appID)
			if err != nil {
				return err
			}
		}
		// Deploy container application
		if serviceType == "container" {
			containerApp := catalog.GetContainerApp()
			err := s.containerService.HandleContainerDeployment(ctx, containerApp, appName, zone, appID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
