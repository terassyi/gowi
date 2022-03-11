package cmd

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/terassyi/gowi/decoder"
)

var dumpCommand = &cobra.Command{
	Use:   "dump",
	Short: "dump WASM file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		d, err := decoder.New(file)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("WASM file: %s\n\n", file)
		r, err := cmd.Flags().GetBool("raw")
		if err != nil {
			log.Fatalln(err)
		}
		if r {
			b, err := decoder.HexDump(file)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(hex.Dump(b))
		}
		s, err := cmd.Flags().GetBool("section")
		if err != nil {
			log.Fatalln(err)
		}
		x, err := cmd.Flags().GetBool("detail")
		if err != nil {
			log.Fatalln(err)
		}
		if s && !x {
			fmt.Println(d.DumpSection())
		}
		if x {
			detail, err := d.DumpDetail()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(detail)
		}
	},
}
