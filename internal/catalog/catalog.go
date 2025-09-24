package catalog

import (
	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/store/model"
)

func GetCatalogVm(name string, serviceName server.ApplicationService) *model.RequestedVm {
	if serviceName == "webserver" {
		return &model.RequestedVm{
			Name:     name,
			Env:      "prod",
			Role:     "webserver",
			TenantId: "tenant-123",
			Ram:      1,
			Cpu:      1,
			Os:       "fedora",
			Region:   "us-east-1",
		}
	} else {
		return &model.RequestedVm{}
	}
}
