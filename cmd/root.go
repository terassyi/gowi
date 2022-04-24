package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gowi [subcommand]",
	Short: "Gowi is a WebAssembly interpreter",
}

var (
	file     string
	sections bool
)

func init() {
	// dump subcommand
	dumpCommand.Flags().BoolP("section", "s", false, "Show sections in WASM file.")
	dumpCommand.Flags().BoolP("raw", "r", false, "Show raw binary.")
	dumpCommand.Flags().BoolP("detail", "x", false, "Show section details.")
	rootCmd.AddCommand(dumpCommand)
	// exec subcommand
	execCommand.Flags().BoolP("list-all-exports", "l", false, "Show all exported functions.")
	execCommand.Flags().StringP("invoke", "i", "", "Invoke an exported function.")
	execCommand.Flags().IntP("debug", "d", 0, "Debug the invoked function.")
	execCommand.Flags().StringSliceP("args", "a", []string{}, "Arguments for the invoking function.")
	rootCmd.AddCommand(execCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
