package v1alpha1

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/service"
)

type ServiceHandler struct {
	srv *service.HealthService
}

func NewServiceHandler(sourceService *service.HealthService) *ServiceHandler {
	return &ServiceHandler{
		srv: sourceService,
	}
}

// (GET /health)
func (s *ServiceHandler) Health(ctx context.Context, request server.HealthRequestObject) (server.HealthResponseObject, error) {
	_ = s.srv.Health(ctx)
	return server.Health200Response{}, nil
}
