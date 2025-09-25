package tier2

default inputvalid := false
default outputvalid := false

result contains r if {
	r := {
		"input": input,
		"output": object.union(input, {"zones": data.t2.zones}),
	}
}

outputvalid if {
	result[r]

	"us-west-1" in r.output.zones
	"us-west-2" in r.output.zones
}

inputvalid if {
	result[r]

	"us-west-1" in r.input.zones
	"us-west-2" in r.input.zones
}
