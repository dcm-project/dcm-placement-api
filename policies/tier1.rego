package tier1

default inputvalid := false

default outputvalid := false

result contains r if {
	input.tier == "1"
	r := {
		"input": input,
		"output": object.union(input, data.env.prod),
	}
}

result contains r if {
	input.tier != "1"
	r := {
		"input": input,
		"output": object.union(input, data.env.dev),
	}
}

# Tier1 validation:
# Ensure the VM is on two specific zones
outputvalid if {
	result[r]

	"us-east-1" in r.output.zones
	"us-east-2" in r.output.zones
}

outputvalid if {
	result[r]

	"us-west-0" in r.output.zones
	"us-west-1" in r.output.zones
}

# Ensure the VM is on two zones
inputvalid if {
	result[r]

	"us-east-1" in r.input.zones
	"us-east-2" in r.input.zones
}

inputvalid if {
	result[r]

	"us-west-0" in r.input.zones
	"us-west-1" in r.input.zones
}
