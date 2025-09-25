package catalog

import (
	"github.com/dcm-project/dcm-placement-api/internal/api/server"
)

type CatalogVm struct {
	Ram int
	Cpu int
	Os  string
}

func GetCatalogVm(serviceName server.ApplicationService) *CatalogVm {
	if serviceName == "webserver" {
		return &CatalogVm{
			Ram: 1,
			Cpu: 1,
			Os:  "fedora",
		}
	}

	return nil
}
