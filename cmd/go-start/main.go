package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set by build flags
	Version = "dev"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-start",
		Short: "A Go project scaffold tool",
		Long: `go-start is a CLI tool that helps you quickly create Go projects
with MVC or DDD architecture, powered by Gin framework.`,
		Version: Version,
	}

	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newSpecCmd())
	rootCmd.AddCommand(newGenCmd())
	rootCmd.AddCommand(newDoctorCmd())
	rootCmd.AddCommand(newVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("go-start version %s\n", Version)
		},
	}
}
