package subnet

restricted_subnets = {"10.1.0.0/24", "10.0.1.0/24"}

default allow = false

allow if {
    input.env == "PROD"
	not restricted(input.network)
}

restricted(ip) if {
    subnet := restricted_subnets[_]
    net.cidr_contains(subnet, ip)
}
