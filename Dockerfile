ARG GOLANG_VERSION=1.25.5

ARG TARGETOS
ARG TARGETARCH

ARG COMMIT
ARG VERSION

FROM --platform=${TARGETARCH} docker.io/golang:${GOLANG_VERSION} AS build

WORKDIR /prometheus-mcp-server

COPY go.* ./
COPY cmd/server ./cmd/server
COPY config ./config
COPY errors ./errors
COPY handlers ./handlers
COPY management ./management
COPY testdata ./testdata

ARG TARGETOS
ARG TARGETARCH

ARG VERSION
ARG COMMIT

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
    -ldflags "-X main.OSVersion=${VERSION} -X main.GitCommit=${COMMIT}" \
    -a -installsuffix cgo \
    -o /go/bin/server \
    ./cmd/server

FROM --platform=${TARGETARCH} gcr.io/distroless/static-debian12:latest

LABEL org.opencontainers.image.description="Prometheus MCP server`"
LABEL org.opencontainers.image.source="https://github.com/DazWilkin/prometheus-mcp-server"

COPY --from=build /go/bin/server /

EXPOSE 7777 8080

ENTRYPOINT ["/server"]
CMD ["--server.addr=:7777","--metric.addr=:8080","--prometheus=http://localhost:9090"]
