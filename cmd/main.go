package main

import (
	"ca-generator/internal/structure_gen"

	"github.com/sirupsen/logrus"
)

// test
func main() {
	generator, err := structure_gen.NewGenerator("settings.json")
	if err != nil {
		logrus.Fatalf("Error creating generator: %v", err)
	}
	generator.Generate()
}
