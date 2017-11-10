package cmd

import(
	"github.com/spf13/cobra"
	"github.com/idahoakl/go-i2c"
	"log"
	"fmt"
)

func init() {
	RootCmd.AddCommand(readCmd)
}

var readCmd = &cobra.Command{
	Use: "read",
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		bus := parseInt(args[0])
		addr := parseInt(args[1])
		readCount := parseInt(args[2])

		if i2C, e := i2c.NewI2C(bus); e != nil {
			log.Fatal(e)
		} else {
			data := make([]byte, readCount)
			if _, e := i2C.Read(uint8(addr), data); e != nil {
				log.Fatal(e)
			} else {
				fmt.Printf("0x%X\n", data)
			}
		}
	},
}