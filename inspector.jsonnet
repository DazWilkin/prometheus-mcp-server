local name = std.extVar("NAME");
local token = std.extVar("TOKEN");
local tailnet = std.extVar("TAILNET");

local labels = {
  "app": name,
  "system": "mcp",
  "type": "server",
  "component": "inspector",
};

local image = "ghcr.io/modelcontextprotocol/inspector:0.16.2";

// Converts a host name into a fully-qualified Tailnet (Ingress) URL
local fqdn(host, tailnet) = "https://%(host)s.%(tailnet)s" % {
  "host": host,
  "tailnet": tailnet,
};

// Defines the base configuration for the web UI and proxy
// Each key will be extended with a unique host and FQDN
local base_config = {
  "webui": {
    "port": 6274,
  },
  "proxy": {
    "port": 6277,
  },
};

// Enrich the base configuration with host name
// Overcomes challenge in referencing the key value to construct the host value
local config = {
  // Use the base configuration
  // And extend it with the host name
  [key]: base_config[key] {
    "host": "%(name)s-%(key)s" % { "name": name, "key": key },
  }
  for key in std.objectFields(base_config)
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
        "containers": [
          {
            "name": name,
            "image": image,
            "args": [],
            "env": [
              {
                "name": "HOST",
                "value": "0.0.0.0",
              },
              {
                "name": "ALLOWED_ORIGINS",
                "value": fqdn(config.webui.host, tailnet),
              },
              {
                "name": "MCP_AUTO_OPEN_ENABLED",
                "value": "false",
              },
              {
                "name": "MCP_PROXY_AUTH_TOKEN",
                "value": token,
              },
            ],
            "ports": [
              {
                "name": variant,
                "containerPort": config[variant].port,
                "protocol": "TCP",
              }
              for variant in std.objectFields(config)
            ],
          },
        ],
      },
    },
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
        "name": variant,
        "port": config[variant].port,
        "targetPort": config[variant].port,
        "protocol": "TCP",
      }
      for variant in std.objectFields(config)
    ],
  },
};

// Generate an Ingress for each service port
// Trying to create one Ingress with different paths didn't work
local ingresses = [
  {
    "apiVersion": "networking.k8s.io/v1",
    "kind": "Ingress",
    "metadata": {
      "name": config[variant].host,
      "labels": labels { "component": variant },
    },
    "spec": {
      "defaultBackend": {
        "service": {
          // Use the singular service name
          "name": name,
          "port": {
            "number": config[variant].port,
          },
        },
      },
      "ingressClassName": "tailscale",
      "tls": [
        {
          "hosts": [
            config[variant].host,
          ],
        },
      ],
    },
  }
  for variant in std.objectFields(config)
];

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
    "name": name,
    "labels": labels,
  },
  "items": [
    deployment,
    service,
    vpa,
  ] + ingresses,
}
