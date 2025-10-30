package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dcm-project/dcm-placement-api/internal/catalog"
	"go.uber.org/zap"
)

type Service struct {
	client *ClientWithResponses
	logger *zap.SugaredLogger
}

func NewService(baseURL string) (*Service, error) {
	client, err := NewClientWithResponses(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider client: %w", err)
	}

	return &Service{
		client: client,
		logger: zap.S().Named("provider_service"),
	}, nil
}

// CreateVMDeployment creates a VM deployment in the provider service
func (s *Service) CreateVMDeployment(ctx context.Context, name, namespace string, vm *catalog.CatalogVm, appID string) (string, error) {
	s.logger.Infow("Creating VM deployment", "name", name, "namespace", namespace)

	// Build the deployment request
	kind := DeploymentRequestKindVm
	labels := map[string]string{
		"app-id": appID,
	}

	vmSpec := VMSpec{
		Vm: struct {
			Cpu int        `json:"cpu"`
			Os  VMSpecVmOs `json:"os"`
			Ram int        `json:"ram"`
		}{
			Ram: vm.Ram,
			Cpu: vm.Cpu,
			Os:  VMSpecVmOs(vm.Os),
		},
	}

	var spec DeploymentRequest_Spec
	if err := spec.FromVMSpec(vmSpec); err != nil {
		return "", fmt.Errorf("failed to create spec from VMSpec: %w", err)
	}

	req := CreateDeploymentJSONRequestBody{
		Kind: kind,
		Metadata: Metadata{
			Name:      name,
			Namespace: &namespace,
			Labels:    &labels,
		},
		Spec: spec,
	}

	// Call the provider service
	resp, err := s.client.CreateDeploymentWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create VM deployment: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		if resp.JSON400 != nil {
			return "", fmt.Errorf("bad request: %s - %s", resp.JSON400.Code, resp.JSON400.Message)
		}
		if resp.JSON409 != nil {
			return "", fmt.Errorf("deployment already exists: %s - %s", resp.JSON409.Code, resp.JSON409.Message)
		}
		if resp.JSON500 != nil {
			return "", fmt.Errorf("internal server error: %s - %s", resp.JSON500.Code, resp.JSON500.Message)
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	if resp.JSON201 == nil || resp.JSON201.Id == nil {
		return "", fmt.Errorf("deployment created but no ID returned")
	}

	s.logger.Infow("VM deployment created successfully", "deploymentID", *resp.JSON201.Id)
	return *resp.JSON201.Id, nil
}

// CreateContainerDeployment creates a container deployment in the provider service
func (s *Service) CreateContainerDeployment(ctx context.Context, name, namespace string, app *catalog.ContainerApp, appID string) (string, error) {
	s.logger.Infow("Creating container deployment", "name", name, "namespace", namespace)

	// Build the deployment request
	kind := DeploymentRequestKindContainer
	labels := map[string]string{
		"app-id": appID,
	}

	replicas := int(app.Replica)
	containerPort := app.Port
	servicePort := app.Port
	protocol := TCP

	containerSpec := ContainerSpec{
		Container: struct {
			Environment *[]struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"environment,omitempty"`
			Image string `json:"image"`
			Ports *[]struct {
				ContainerPort int                                  `json:"containerPort"`
				Protocol      *ContainerSpecContainerPortsProtocol `json:"protocol,omitempty"`
				ServicePort   *int                                 `json:"servicePort,omitempty"`
			} `json:"ports,omitempty"`
			Replicas  *int `json:"replicas,omitempty"`
			Resources *struct {
				Cpu    *string `json:"cpu,omitempty"`
				Memory *string `json:"memory,omitempty"`
			} `json:"resources,omitempty"`
		}{
			Image:    app.Image,
			Replicas: &replicas,
			Ports: &[]struct {
				ContainerPort int                                  `json:"containerPort"`
				Protocol      *ContainerSpecContainerPortsProtocol `json:"protocol,omitempty"`
				ServicePort   *int                                 `json:"servicePort,omitempty"`
			}{
				{
					ContainerPort: containerPort,
					ServicePort:   &servicePort,
					Protocol:      &protocol,
				},
			},
		},
	}

	var spec DeploymentRequest_Spec
	if err := spec.FromContainerSpec(containerSpec); err != nil {
		return "", fmt.Errorf("failed to create spec from ContainerSpec: %w", err)
	}

	req := CreateDeploymentJSONRequestBody{
		Kind: kind,
		Metadata: Metadata{
			Name:      name,
			Namespace: &namespace,
			Labels:    &labels,
		},
		Spec: spec,
	}

	// Call the provider service
	resp, err := s.client.CreateDeploymentWithResponse(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create container deployment: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		if resp.JSON400 != nil {
			return "", fmt.Errorf("bad request: %s - %s", resp.JSON400.Code, resp.JSON400.Message)
		}
		if resp.JSON409 != nil {
			return "", fmt.Errorf("deployment already exists: %s - %s", resp.JSON409.Code, resp.JSON409.Message)
		}
		if resp.JSON500 != nil {
			return "", fmt.Errorf("internal server error: %s - %s", resp.JSON500.Code, resp.JSON500.Message)
		}
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	if resp.JSON201 == nil || resp.JSON201.Id == nil {
		return "", fmt.Errorf("deployment created but no ID returned")
	}

	s.logger.Infow("Container deployment created successfully", "deploymentID", *resp.JSON201.Id)
	return *resp.JSON201.Id, nil
}

// DeleteDeployment deletes a deployment by ID
func (s *Service) DeleteDeployment(ctx context.Context, deploymentID string) error {
	s.logger.Infow("Deleting deployment", "deploymentID", deploymentID)

	resp, err := s.client.DeleteDeploymentWithResponse(ctx, deploymentID)
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	if resp.StatusCode() != http.StatusNoContent {
		if resp.JSON404 != nil {
			return fmt.Errorf("deployment not found: %s - %s", resp.JSON404.Code, resp.JSON404.Message)
		}
		if resp.JSON500 != nil {
			return fmt.Errorf("internal server error: %s - %s", resp.JSON500.Code, resp.JSON500.Message)
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	s.logger.Infow("Deployment deleted successfully", "deploymentID", deploymentID)
	return nil
}
