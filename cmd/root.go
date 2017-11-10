package cmd

import(
	"github.com/spf13/cobra"
	"strconv"
	"log"
)

var RootCmd = &cobra.Command{
	Use: "i2c-test",
}

func parseInt(s string) int {
	if i, e := strconv.ParseInt(s, 0, 8); e != nil {
		log.Fatal(e)
		return -1
	} else {
		return int(i)
	}
}