package mappers

import (
	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/store/model"
)

func ApplicationListToAPI(dbApps model.ApplicationList) server.ApplicationList {
	var apiApps server.ApplicationList
	for _, dbApp := range dbApps {
		zones := []string(dbApp.Zones)
		apiApp := server.Application{
			Name:    dbApp.Name,
			Service: server.ApplicationService(dbApp.Service),
			Tier:    &dbApp.Tier,
			Zones:   &zones,
		}
		apiApps = append(apiApps, apiApp)
	}
	return apiApps
}
