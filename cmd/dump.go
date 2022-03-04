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
		d := decoder.New(file)
		m, err := d.Decode()
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
			fmt.Println(m.Dump())
		}
		if x {
			d, err := m.DumpDetail()
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(d)
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
