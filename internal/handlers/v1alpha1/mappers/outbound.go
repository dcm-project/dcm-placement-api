package mappers

import (
	"github.com/dcm-project/dcm-placement-api/internal/api/server"
	"github.com/dcm-project/dcm-placement-api/internal/store/model"
)

// RequestedVmToAPI converts database RequestedVm to API VM
func RequestedVmToAPI(dbVm model.RequestedVm) server.VM {
	return server.VM{
		Name:     dbVm.Name,
		Env:      dbVm.Env,
		Ram:      dbVm.Ram,
		Os:       dbVm.Os,
		Cpu:      dbVm.Cpu,
		Region:   dbVm.Region,
		Role:     dbVm.Role,
		TenantId: &dbVm.TenantId,
	}
}

// RequestedVmListToAPI converts database RequestedVmList to API RequestedVmList
func RequestedVmListToAPI(dbVms model.RequestedVmList) server.RequestedVmList {
	var apiVms server.RequestedVmList
	for _, dbVm := range dbVms {
		apiVm := RequestedVmToAPI(dbVm)
		apiVms = append(apiVms, apiVm)
	}
	return apiVms
}

// DeclaredVmToAPI converts database DeclaredVm to API DeclaredVm
func DeclaredVmToAPI(dbVm model.DeclaredVm) server.DeclaredVm {
	idStr := dbVm.ID.String()
	createdAt := dbVm.CreatedAt.Format("2006-01-02T15:04:05")
	result := server.DeclaredVm{
		Id:        &idStr,
		IpAddress: &dbVm.IPAddress,
		Gateway:   &dbVm.Gateway,
		Netmask:   &dbVm.Netmask,
		DnsName:   &dbVm.DnsName,
		CreatedAt: &createdAt,
		// Include preloaded RequestedVm data
		Name:     dbVm.RequestedVm.Name,
		Env:      dbVm.RequestedVm.Env,
		Ram:      dbVm.RequestedVm.Ram,
		Os:       dbVm.RequestedVm.Os,
		Cpu:      dbVm.RequestedVm.Cpu,
		Region:   dbVm.RequestedVm.Region,
		Role:     dbVm.RequestedVm.Role,
		TenantId: &dbVm.RequestedVm.TenantId,
	}

	return result
}

// DeclaredVmListToAPI converts database DeclaredVmList to API DeclaredVmList
func DeclaredVmListToAPI(dbVms model.DeclaredVmList) server.DeclaredVmList {
	var apiVms server.DeclaredVmList
	for _, dbVm := range dbVms {
		apiVm := DeclaredVmToAPI(dbVm)
		apiVms = append(apiVms, apiVm)
	}
	return apiVms
}
