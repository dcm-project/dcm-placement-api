package v1alpha1

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/service"
	"go.uber.org/zap"
)

type ServiceHandler struct {
	srv *service.HealthService
	ps  *service.PlacementService
}

func NewServiceHandler(
	sourceService *service.HealthService,
	placementService *service.PlacementService) *ServiceHandler {
	return &ServiceHandler{
		srv: sourceService,
		ps:  placementService,
	}
}

// (GET /health)
func (s *ServiceHandler) Health(ctx context.Context, request server.HealthRequestObject) (server.HealthResponseObject, error) {
	_ = s.srv.Health(ctx)
	return server.Health200Response{}, nil
}

// (POST /place/vm)
func (s *ServiceHandler) PlaceVM(ctx context.Context, request server.PlaceVMRequestObject) (server.PlaceVMResponseObject, error) {
	logger := zap.S().Named("placement_handler")
	logger.Info("Processing VM Placement Request")

	err := s.ps.PlaceVM(ctx, request.Body)
	if err != nil {
		logger.Error("Failed to Place Virtual Machine")
		return server.PlaceVM400JSONResponse{}, nil
	}
	logger.Info("Successfully Place Virtual Machine")
	return server.PlaceVM200JSONResponse{}, nil
}
