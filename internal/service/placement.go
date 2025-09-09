package service

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
)

type PlacementService struct {
}

func NewPlacementService() *PlacementService {
	return &PlacementService{}
}

func (s *PlacementService) PlaceVM(ctx context.Context, request *server.PlaceVMRequestObject) error {
	_ = request
	// TODO: validate request with opa
	// TODO: store in db if successful validation
	return nil
}
