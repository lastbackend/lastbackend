package main

import (
	"github.com/lastbackend/registry/pkg/registry"
)

type Analyzer struct {
	Kind string `yaml:"kind"`
}

func main() {
	registry.Run()
}
