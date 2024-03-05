package main

import "github.com/Pii-Hole-Engineering/no-phi-ai/pkg/manager"

// main() function for no-phi-ai app is minimal by design
func main() {
	// setup a new Manager for the app
	m := manager.New()
	// initialize the manager based on the configuration
	m.Init()
	// run the app in the configured mode
	m.Run()
}
