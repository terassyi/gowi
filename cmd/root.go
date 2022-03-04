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
	dumpCommand.Flags().BoolP("section", "s", false, "Show sections in WASM file.")
	dumpCommand.Flags().BoolP("raw", "r", false, "Show raw binary.")
	dumpCommand.Flags().BoolP("detail", "x", false, "Show section details.")
	rootCmd.AddCommand(dumpCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
