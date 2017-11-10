package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/idahoakl/go-i2c"
	"encoding/hex"
)

func init() {
	RootCmd.AddCommand(writeCmd)
}

var writeCmd = &cobra.Command{
	Use: "write",
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		bus := parseInt(args[0])
		addr := parseInt(args[1])
		writeData := parseHex(args[2])

		if i2C, e := i2c.NewI2C(bus); e != nil {
			log.Fatal(e)
		} else {
			if _, e := i2C.Write(uint8(addr), writeData); e != nil {
				log.Fatal(e)
			}
		}
	},
}

func parseHex(s string) []byte {
	if b, e := hex.DecodeString(s); e != nil {
		log.Fatal(e)
		return nil
	} else {
		return b
	}
}