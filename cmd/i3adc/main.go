package main

import (
	"github.com/seeruk/i3adc/internal"
)

func main() {
	resolver := internal.NewResolver()

	logger := resolver.ResolveLogger()
	logger.Info("main: i3adc starting...")

	// ...

	logger.Info("main: i3adc exiting...")
}
