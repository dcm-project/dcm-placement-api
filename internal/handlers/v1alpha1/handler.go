package v1alpha1

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/handlers/v1alpha1/mappers"
	"github.com/dcm-project/dcm-placement-api/internal/service"
	"github.com/dcm-project/dcm-placement-api/internal/store"
	"go.uber.org/zap"
)

type ServiceHandler struct {
	ps    *service.PlacementService
	store store.Store
}

func NewServiceHandler(store store.Store, placementService *service.PlacementService) *ServiceHandler {
	return &ServiceHandler{
		store: store,
		ps:    placementService,
	}
}

// (GET /health)
func (s *ServiceHandler) Health(ctx context.Context, request server.HealthRequestObject) (server.HealthResponseObject, error) {
	return server.Health200Response{}, nil
}

// (GET /applications)
func (s *ServiceHandler) GetApplications(ctx context.Context, request server.GetApplicationsRequestObject) (server.GetApplicationsResponseObject, error) {
	applications, err := s.store.Application().List(ctx)
	if err != nil {
		return server.GetApplications400JSONResponse{}, err
	}
	return server.GetApplications200JSONResponse(mappers.ApplicationListToAPI(applications)), nil
}

// (POST /applications)
func (s *ServiceHandler) CreateApplication(ctx context.Context, request server.CreateApplicationRequestObject) (server.CreateApplicationResponseObject, error) {
	logger := zap.S().Named("placement_service")
	logger.Info("Creating Application. ", "Application: ", request)

	app, err := s.ps.CreateApplication(ctx, request.Body)
	if err != nil {
		logger.Error("Failed to create Application: ", "error", err)
		return server.CreateApplication400JSONResponse{}, err
	}
	logger.Info("Application created. ", "Application: ", app)
	return server.CreateApplication201JSONResponse(*app), nil
}
