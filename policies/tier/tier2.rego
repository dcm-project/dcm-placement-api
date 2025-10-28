package tier2

import rego.v1

default valid := false

# Helper rule to get the server URL from data with fallback
server_url := data.server_url if {
    data.server_url
}

# Default server URL if not provided in data
default server_url := "http://localhost:8080"

# Helper function to fetch zones for a given set of labels
fetch_zones(labels) := zone_names if {
    response := http.send({
        "method": "POST",
        "url": sprintf("%s/namespaces", [server_url]),
        "headers": {
            "Content-Type": "application/json"
        },
        "body": {
            "labels": labels
        }
    })

    # Check if the API call was successful
    response.status_code == 200

    # Extract zone names from the response
    zone_names := [ns.name | ns := response.body.namespaces[_]]
}

# Default return empty list if fetch fails
default fetch_zones(_) := []

# Required zones for tier 2 - try production labels first
required_zones := zone_names if {
    production_labels := data.t2.production_labels
    production_zones := fetch_zones(production_labels)
    count(production_zones) > 0
    zone_names := production_zones
}

# Fallback to backup labels if production returns empty
required_zones := zone_names if {
    production_labels := data.t2.production_labels
    production_zones := fetch_zones(production_labels)
    count(production_zones) == 0

    backup_labels := data.t2.backup_labels
    zone_names := fetch_zones(backup_labels)
}

# Generate failures if zones are defined but not equal to required_zones
failures contains failure if {
    input.zones  # Zones field exists
    some zone in required_zones
    not zone in input.zones
    failure := sprintf("Missing required zone '%s' in input specification", [zone])
}

failures contains failure if {
    input.zones  # Zones field exists
    some zone in input.zones
    not zone in required_zones
    failure := sprintf("Unexpected zone '%s' in input specification", [zone])
}

# Input is valid if zones are not defined OR zones exactly match required_zones
valid if {
    not input.zones  # Zones field does not exist - this is valid
}

valid if {
    input.zones  # Zones field exists
    # All required zones are present
    every zone in required_zones {
        zone in input.zones
    }
    # No extra zones are present
    every zone in input.zones {
        zone in required_zones
    }
}
