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

// (GET /requestedvms)
func (s *ServiceHandler) GetRequestedVms(ctx context.Context, request server.GetRequestedVmsRequestObject) (server.GetRequestedVmsResponseObject, error) {
	vms, err := s.store.RequestedVm().List(ctx)
	if err != nil {
		return server.GetRequestedVms400JSONResponse{}, err
	}

	return server.GetRequestedVms200JSONResponse(mappers.RequestedVmListToAPI(vms)), nil
}

// (GET /declaredvms)
func (s *ServiceHandler) GetDeclaredVms(ctx context.Context, request server.GetDeclaredVmsRequestObject) (server.GetDeclaredVmsResponseObject, error) {
	vms, err := s.store.DeclaredVm().List(ctx)
	if err != nil {
		return server.GetDeclaredVms400JSONResponse{}, err
	}

	return server.GetDeclaredVms200JSONResponse(mappers.DeclaredVmListToAPI(vms)), nil
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
	logger.Info("Creating Application", "Application", request)

	err := s.ps.CreateApplication(ctx, request.Body)
	if err != nil {
		logger.Error("Failed to create Application: ", "error", err)
		return server.CreateApplication400JSONResponse{}, err
	}
	logger.Info("Application created", "Application", request)
	return server.CreateApplication201JSONResponse{
		Name:    request.Body.Name,
		Service: request.Body.Service,
	}, nil
}

// (POST /place/vm)
func (s *ServiceHandler) PlaceVM(ctx context.Context, request server.PlaceVMRequestObject) (server.PlaceVMResponseObject, error) {
	logger := zap.S().Named("placement_handler")
	logger.Info("Processing VM Placement Request")

	err := s.ps.PlaceVM(ctx, request.Body)
	if err != nil {
		logger.Error("Failed to place Virtual Machine: ", "error", err)
		return server.PlaceVM400JSONResponse{Error: err.Error()}, nil
	}
	logger.Info("Successfully Place Virtual Machine")
	return server.PlaceVM200JSONResponse{Message: "VM successfully placed"}, nil

}
