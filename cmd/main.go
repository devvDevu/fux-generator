package main

import (
	"ca-generator/internal/structure_gen"
	"time"

	"github.com/sirupsen/logrus"
)

// test
func main() {
	t := time.Now()
	generator, err := structure_gen.NewGenerator("settings.json")
	if err != nil {
		logrus.Fatalf("Error creating generator: %v", err)
	}
	generator.Generate()
	timeEnd := time.Since(t)
	logrus.Infof("Time: %v", timeEnd)
}
