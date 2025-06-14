package main

import (
	"flag"
	"fmt"
)

type Config struct {
	Prometheus string
	Metrics    Metrics
}

func NewConfig() *Config {
	addr := flag.String("metrics.addr", ":8080", "Endpoint on which metrics are published")
	path := flag.String("metrics.path", "/metrics", "Path on which metrics are served")
	prometheus := flag.String("prometheus", "http://localhost:9090", "Endpoint of Prometheus server")
	flag.Parse()

	return &Config{
		Prometheus: *prometheus,
		Metrics: Metrics{
			Addr: *addr,
			Path: *path,
		},
	}
}

type Metrics struct {
	Addr string
	Path string
}

func (m Metrics) GoString() string {
	return fmt.Sprintf("Metrics{Endpoint: %q, Path: %q}", m.Addr, m.Path)
}
func (m Metrics) String() string {
	return fmt.Sprintf("%s/%s", m.Addr, m.Path)
}
