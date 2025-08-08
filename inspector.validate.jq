# inspector.validate.jq

# Expects variables from the invocation:
# jq \
# --argjson TLS [true|false]
# --arg NAME [string]
# --arg TAILNET [string]

# These are piped into the filter
6274 as $WEBUI_PORT |
6277 as $PROXY_PORT |

# Each of these is a test
# Each test is represented by an
# if
#   {predicate}
# then output the result
# else throw an error

# Debugging
# Replace predicates or branches with e.g. true
# if ( true) then . else error("message") end
# if ( true | true | ( true) ) then . else error("message") end

# Items list check
if (
    # Assert that the output contains items
    .items
    | type=="array"
)
then .
else error("Items list check")
end

and

# Deployment check
if (
    # Assert that there is a Deployment kind
    .items
    | any(.kind=="Deployment")
)
then .
else error("Deployment check")
end

and

# ALLOWED_ORIGINS check
if (
    # Assert that the Deployment has correct ALLOWED_ORIGINS
    .items[]
    | select(.kind=="Deployment" and .metadata.name==$NAME)
    | .spec.template.spec.containers[]
    | select(.name==$NAME)
    | .env[]
    | select(.name=="ALLOWED_ORIGINS")
    | (
        if $TLS
        then
            "https://\($NAME)-webui.\($TAILNET)"
        else
            "http://\($NAME).\($TAILNET):\($WEBUI_PORT)"
        end
      ) as $VALUE
    | .value==$VALUE
) then . else error("ALLOWED_ORIGINS check") end

and

# Service check
if (
    # Assert that there is a Service kind
    .items | any(.kind=="Service")
)
then .
else error("Service check")
end

and

# Assert that the Service is correctly configured if NOT TLS
if (
    $TLS==false
)
then (
    if (
    .items[]
    | select(.kind=="Service" and .metadata.name==$NAME)
    | (
        .metadata.annotations["tailscale.com/hostname"]==$NAME
        ) and (
        .spec.type=="LoadBalancer"
        ) and (
        .spec.loadBalancerClass=="tailscale"
        )
    )
    then .
    else error("Service check")
    end
)
else .
end

and

# Ingress check
if (
    # Assert that if TLS then there are 2 Ingress, 0 otherwise
    .items
    | map(select(.kind=="Ingress"))
    | length
    | if $TLS then .==2 else .==0 end
)
then .
else error("Ingress check")
end

# Pipe the results into valid (true) or invalid (false)
| if . == true then "valid" else "invalid" end