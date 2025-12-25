package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the project in development mode",
		Long: `Run the project with hot reload support using air.

If air is not installed, it will run the project directly without hot reload.`,
		RunE: runRun,
	}

	return cmd
}

func runRun(cmd *cobra.Command, args []string) error {
	// Check if we're in a Go project
	if !isGoProject() {
		return fmt.Errorf("not in a Go project directory")
	}

	// Check if air is installed for hot reload
	if hasCommand("air") {
		fmt.Println("Running with hot reload (air)...")
		return runWithAir()
	}

	// Run directly
	fmt.Println("Running without hot reload (install air for hot reload support: go install github.com/cosmtrek/air@latest)")
	return runDirectly()
}

func runWithAir() error {
	cmd := exec.Command("air")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runDirectly() error {
	// Try to find main.go
	mainPaths := []string{
		"cmd/server/main.go",
		"main.go",
	}

	var mainPath string
	for _, path := range mainPaths {
		if _, err := os.Stat(path); err == nil {
			mainPath = path
			break
		}
	}

	if mainPath == "" {
		return fmt.Errorf("main.go not found (tried: %v)", mainPaths)
	}

	cmd := exec.Command("go", "run", mainPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func isGoProject() bool {
	if _, err := os.Stat("go.mod"); err != nil {
		return false
	}
	return true
}

func hasCommand(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.Success()
}
