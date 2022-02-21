package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terassyi/gowi/decoder"
)

var dumpCommand = &cobra.Command{
	Use:   "dump",
	Short: "dump WASM file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		r, err := cmd.Flags().GetBool("section")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := dump(file, r); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func dump(file string, sections bool) error {
	d := decoder.New(file)
	_, err := d.Decode()
	if err != nil {
		return err
	}
	return nil
}
