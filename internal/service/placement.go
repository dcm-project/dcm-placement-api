package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/vm_subnet"
	"go.uber.org/zap"
)

type PlacementService struct {
}

func NewPlacementService() *PlacementService {
	return &PlacementService{}
}

func (s *PlacementService) PlaceVM(ctx context.Context, request *server.PlaceVMJSONRequestBody) error {
	logger := zap.S().Named("placement_service")
	logger.Info("Processing Placement request", "VM-NAME", request.Name)
	// customize network based on conditions
	subnets := vm_subnet.GetSubnetList()
	logger.Info("Subnet List", "Subnets", subnets)

	spec := new(vm_subnet.Spec)
	for _, s := range subnets {
		conditions := s.VMConditions
		if (conditions.Role == request.Role) &&
			(conditions.Region == request.Region) &&
			(conditions.Environment == request.Env) &&
			(conditions.TenantId == *request.TenantId) {

			networkSpec := s.NetworkSpec
			spec.DnsName = networkSpec.DnsName
			spec.Gateway = networkSpec.Gateway
			spec.IPAddress = networkSpec.IPAddress
			spec.Netmask = networkSpec.Netmask

			break
		}
	}
	if reflect.DeepEqual(spec, new(vm_subnet.Spec)) {
		logger.Error("Failed VM placement - Subnet conditions not met")
		return fmt.Errorf("failed to meet any subnet conditions")
	}
	logger.Info("Processed network spec for vm place request", "VM", request, "Network-Spec", spec)
	// TODO: validate request with opa
	// TODO: store in db if successful validation
	return nil
}
