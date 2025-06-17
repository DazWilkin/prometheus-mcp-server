local name = std.extVar("NAME");

local labels = {
    "app": name,
    "system": "mcp",
    "type": "server",
    "upstream": "prometheus",
};

// Needs to remain as-is so that GitHub Workflow can replace it on updates
local image = "ghcr.io/dazwilkin/prometheus-mcp-server:3209e69b2dfa5b3198044cd94380123a11b6d656";

// port must be a number
// So that it can be referenced as such in Deployment|Service specs etc.
local server(host, port) = {
    "host": host,
    "port": port,
};

// Expects a server object with host (string) and port (number) properties
// Returns a string in the format "host:port"
local addr(server) = "%(host)s:%(port)d" % server;

// Represents server configurations:
// 1. MCP server
// 2. Metrics server
// 3. Prometheus URL
local config = {
    "server": server("0.0.0.0", 7777),
    "metric": server("0.0.0.0", 8080),
    "prometheus": "http://prometheus-operated.monitoring.svc.cluster.local:9090",
};

// Represents GHCR authentication
local ghcr = {
    local registry = "https://ghcr.io",

    local username = std.extVar("GHCR_USERNAME"),
    local password = std.extVar("GHCR_TOKEN"),
    local email = std.extVar("GHCR_EMAIL"),

    // Value must be base64-encoded
    local auth = std.base64("%(username)s:%(password)s" % {
        "username": username,
        "password": password,
    }),

    "auths": {
        // Must be [registry] to be evaluated correctly
        [registry]: {
            "username": username,
            "password": password,
            "email": email,
            "auth": auth,
        },
    },
};

local deployment = {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "spec": {
        "selector": {
            "matchLabels": labels,
        },
        "template": {
            "metadata": {
                "name": name,
                "labels": labels,
            },
            "spec": {
                "serviceAccount": name,
                "containers": [
                    {
                        "name": name,
                        "image": image,
                        "args": [
                            // Addresses must be Pod local (0.0.0.0)
                            // The respective (external) host addresses are used by scripts
                            "--server.addr=%(addr)s" % addr(config.server),
                            "--metric.addr=%(addr)s" % addr(config.metric),
                            "--prometheus=%(prometheus)s" % { "prometheus": config.prometheus },
                            // Defaults need not be set
                            // "--server.path=/mcp",
                            // "--metric.path="/metrics",
                        ]
                    },
                ],
            },
        },
    },
};

local ingress = {
    "apiVersion": "networking.k8s.io/v1",
    "kind": "Ingress",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "spec": {
        "defaultBackend": {
            "service": {
                "name": name,
                "port": {
                    "number": config.server.port,
                },
            },
        },
        "ingressClassName": "tailscale",
        "tls": [
            {
                "hosts": [
                    name,
                ],
            },
        ],
    },
};

local rule = {
    "apiVersion": "monitoring.coreos.com/v1",
    "kind": "PrometheusRule",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "spec": {
        "groups": [
            {
                "name": name,
                "rules": [
                    {
                        local minutes = 5,
                        "alert": "PrometheusMCPServerDown",
                        "expr": "up{job=\"mcp-server\"} == 0",
                        "for": "%(for)dm" % { "for": minutes },
                        "labels": {
                            "severity": "critical"
                        },
                        "annotations": {
                            "summary": "Prometheus MCP server is down",
                            "description": "Prometheus MCP server has been down for more than 5 minutes."
                        },
                    },
                    {
                        local minutes = 5,
                        "alert": "PrometheusMCPToolErrors",
                        "expr": "mcp_prometheus_error",
                        "for": "%(for)dm" % { "for": minutes },
                        "labels": {
                            "severity": "warning"
                        },
                        "annotations": {
                            "summary": "Prometheus MCP tool reporting errors",
                            "description": "Prometheus MCP tool ({{ $labels.tool }}) reporting errors ({{ $value }})"
                        },
                    },
                ],
            },
        ],
    },
};

local secret = {
    "apiVersion": "v1",
    "kind": "Secret",
    "metadata": {
        "name": name,
    },
    "type": "kubernetes.io/dockerconfigjson",
    "data": {
        ".dockerconfigjson": std.base64(std.manifestJsonEx(ghcr,"")),
    },
};

local service = {
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "spec": {
        "selector": labels,
        "ports": [
            {
                "name": "json-rpc",
                "port": config.server.port,
                "targetPort": config.server.port,
                "protocol": "TCP",
            },
            {
                "name": "metrics",
                "port": config.metric.port,
                "targetPort": config.metric.port,
                "protocol": "TCP",
            }
        ]
    }
};

local service_account = {
    "apiVersion": "v1",
    "kind": "ServiceAccount",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "imagePullSecrets": [
        {
            "name": name,
        },
    ],
};

local service_monitor = {
    "apiVersion": "monitoring.coreos.com/v1",
    "kind": "ServiceMonitor",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "spec": {
        "selector": {
            "matchLabels": labels,
        },
        "endpoints": [
            {
                "path": "/metrics",
                "port": "metrics",
            },
        ],
    },
};

local vpa = {
    "apiVersion": "autoscaling.k8s.io/v1",
    "kind": "VerticalPodAutoscaler",
    "metadata": {
        "name": name,
        "labels": labels,
    },
    "spec": {
        "targetRef": {
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "name": name,
        },
        "updatePolicy": {
            "updateMode": "Off",
        },
    },
};

// Output
{
    "apiVersion": "v1",
    "kind": "List",
    "metadata": {
        "name": "list",
    },
    "items": [
        deployment,
        ingress,
        rule,
        secret,
        service,
        service_account,
        service_monitor,
        vpa,
    ],
}
