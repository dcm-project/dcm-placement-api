package tier1

default inputvalid := false
default outputvalid := false

result contains r if {
	r := {
		"input": input,
		"output": object.union(input, {"zones": data.t1.zones}),
	}
}

outputvalid if {
	result[r]

	"us-east-1" in r.output.zones
	"us-east-2" in r.output.zones
}

# Ensure the VM is on two zones
inputvalid if {
	result[r]

	"us-east-1" in r.input.zones
	"us-east-2" in r.input.zones
}
