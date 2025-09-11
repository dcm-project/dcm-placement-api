package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/deploy"
	"github.com/dcm-project/dcm-placement-api/internal/opa"
	"github.com/dcm-project/dcm-placement-api/internal/store"
	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"github.com/dcm-project/dcm-placement-api/internal/vm_subnet"
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

func (s *PlacementService) PlaceVM(ctx context.Context, request *server.PlaceVMJSONRequestBody) error {
	logger := zap.S().Named("placement_service")
	logger.Info("Processing Placement request", "VM-NAME", request.Name)

	// Store request record in db:
	requestedVm, err := s.store.RequestedVm().Create(ctx, model.RequestedVm{
		ID:       uuid.New(),
		Name:     request.Name,
		Env:      request.Env,
		Ram:      request.Ram,
		Os:       request.Os,
		Cpu:      request.Cpu,
		Region:   request.Region,
		Role:     request.Role,
		TenantId: *request.TenantId,
	})
	if err != nil {
		return err
	}

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

	// OPA validation:
	allow, restrictedSubnets, err := s.opa.EvalPolicy(ctx, "subnet", map[string]string{
		"env":     request.Env,
		"network": spec.IPAddress,
	})
	if err != nil {
		return err
	}

	if !allow {
		return fmt.Errorf(
			"Cannot create VM in requested subnet %v, restricted subnets: %v",
			spec.IPAddress,
			restrictedSubnets,
		)
	}

	// Store declared vm in db:
	_, err = s.store.DeclaredVm().Create(ctx, model.DeclaredVm{
		ID:            uuid.New(),
		RequestedVmID: requestedVm.ID,
		IPAddress:     spec.IPAddress,
		Gateway:       spec.Gateway,
		Netmask:       spec.Netmask,
		DnsName:       spec.DnsName,
	})
	if err != nil {
		return err
	}

	// Deploy the vm
	err = s.deploy.DeployVM(ctx, requestedVm)
	if err != nil {
		return err
	}

	return nil
}
