package vm_subnet

import "go.uber.org/zap"

type Subnet struct {
	VMConditions Condition
	NetworkSpec  Spec
}
type Condition struct {
	Region      string
	Environment string
	Role        string
	TenantId    string
}

type Spec struct {
	Gateway   string
	IPAddress string
	Netmask   string
	DnsName   string
}

var subnets = []Subnet{
	{
		VMConditions: Condition{
			Region:      "us-east-1",
			Environment: "PROD",
			Role:        "public-facing",
			TenantId:    "PRCR-001",
		},
		NetworkSpec: Spec{
			Gateway:   "10.0.0.1",
			IPAddress: "10.0.0.12",
			Netmask:   "255.255.255.0",
			DnsName:   "prcr.app.prod.com",
		},
	},
	{
		VMConditions: Condition{
			Region:      "us-east-2",
			Environment: "STAGE",
			Role:        "internal-facing",
			TenantId:    "STCR-001",
		},
		NetworkSpec: Spec{
			Gateway:   "10.0.0.1",
			IPAddress: "10.0.0.12",
			Netmask:   "255.255.255.0",
			DnsName:   "stcr.app.stage.com",
		},
	},
	{
		VMConditions: Condition{
			Region:      "us-east-1",
			Environment: "DEV",
			Role:        "internal-facing",
			TenantId:    "DVCR-001",
		},
		NetworkSpec: Spec{
			Gateway:   "10.0.0.1",
			IPAddress: "10.0.0.12",
			Netmask:   "255.255.255.0",
			DnsName:   "dvcr.app.dev.com",
		},
	},
}

func GetSubnetList() []Subnet {
	zap.S().Info("Getting Subnet Condition List")
	return subnets
}
