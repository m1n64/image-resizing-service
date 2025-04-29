package main

import (
	"context"
	"image-resizing-service/pkg/di"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	_ = context.Background()

}
