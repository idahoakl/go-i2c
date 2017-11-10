package main

import (
	"fmt"
	"os"
	"github.com/idahoakl/go-i2c/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
