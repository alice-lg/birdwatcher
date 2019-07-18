
Minimal description of required fields

# Envelope 

{
    "api": {
        "version": "string",
        "result_from_cache": "boolean",
        "cache_status": {
            "cached_at": "datetime",
            "date": "datetime",
        }
    }
    "ttl": "datetime",
}


# Status
{
    "api": ...,
    "status": {
        "current_server": "datetime",
        "last_reboot": "datetime",
        "last_reconfig": "datetime",
        "version": "string",
        "message": "string",
        "router_id": "string",
    }
}


# Routes

{
    "api": ...,
    "routes": [
        {
            "age": "datetime",
            "bgp": {
                "as_path": ["int"],
                "communities": [["int"]],
                "ext_communities": [["string"]],
                "large_communities": [["int"]],
                "local_pref": "int",
                "med": "int",
                "origin": "string",
                "next_hop": "string",
            },
            "network": "string",
            "from_protocol": "string",
            "interface": "string",
            "gateway": "string"
            "metric": "int",
            "type": ["string"],
            "primary": "boolean"
        }
    ]
}


# Protocols / Neighbors

{
    "api": ...,
    "protocols": [
        {
            "routes": {
                "imported": "int",
                "filtered": "int",
                "exported": "int",
                "preferred": "int",
            },
            "neighbor_address": string,
            "neighbor_as": int,
            "state": "string",
            "description": "string",
            "state_changed": "datetime",
            "uptime": "datetime",
            "last_error": "string"
        }
    ]
}


