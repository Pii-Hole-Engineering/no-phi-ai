package main

import (
	"os"
	"testing"
)

func TestMainApp(t *testing.T) {
	_ = os.Setenv("NOPHI_CONFIG_PATH", "../config/test.yml.example")
	_ = os.Setenv("NOPHI_COMMAND_RUN", "version")
	main()
}
