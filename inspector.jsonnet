// There are (at least) 2 ways to deploy the Inspector
// Each requires different configurations:
// 1. Ingress (TLS=true)
//  + Publishes HTTPS endpoints `inspector-webui` and `inspector-proxy`
//  + Requires Deployment ALLOWED_ORIGINS of `https://inspector-webui.{tailnet}`
//  + Basic Service
//  + 2x Ingress
//
// 2. Service type: LoadBalancer (TLS=false)
//  + Publishes HTTP endpoint `http://inspector.{tailnet}` with WebUI|Proxy ports
//  + Requires Deployment ALLOWED_ORIGINS of `http://inspector.{tailnet}`
//  + Augment Service w/ Tailscale annotations and LoadBalancer class
//  + No Ingress
local TLS = if std.extVar("TLS")=="T" then true else false;

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
local fqdn(host, tailnet) = "%(scheme)s://%(host)s.%(tailnet)s" % {
  // Depends upon the TLS setting (Ingress or Service deployment)
  "scheme": if TLS then "https" else "http",
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
  [variant]: base_config[variant] {
    "host": "%(name)s-%(variant)s" % {
      "name": name,
      "variant": variant,
    },
  }
  for variant in std.objectFields(base_config)
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
                // Depends upon the TLS setting (Ingress or Service deployment)
                "value": if TLS then
                  // If TLS (Ingress) then we want to permit:
                  // https://inspector-webui.{tailnet}
                  fqdn(config.webui.host, tailnet)
                else
                  // If not TLS (Service) then we want to permit:
                  // http://inspector.{tailnet}
                  fqdn(name, tailnet),
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

// Depends upon the TLS setting (Ingress or Service deployment)
// If TLS (Ingress) then create a Service with its base configuration
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
} + (
  // If not TLS (Service) then add Tailscale hostname annotation and LoadBalancer class
  if !TLS then {
    // Uses metadata+: to MERGE into the existing service metadata
    "metadata"+: {
      "annotations": {
        // Overrride the default Tailscale hostname
        // https://tailscale.com/kb/1445/kubernetes-operator-customization#using-custom-machine-names
        "tailscale.com/hostname": name,
      },
    },
    // Uses spec+: to MERGE into the existing service spec
    "spec"+: {
      "type": "LoadBalancer",
      "loadBalancerClass": "tailscale",
    },
  } else {
  }
);

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

local items = [
    deployment,
    service,
    vpa,
  ]
  + if TLS then ingresses else [];

// Output
{
  "apiVersion": "v1",
  "kind": "List",
  "metadata": {
    "name": name,
    "labels": labels,
  },
  "items": items,
}
