package main

import (
	"fmt"
	"os"

	"github.com/clcng/bitcoin-wallet/cmd"
	"github.com/spf13/viper"
)

var (
	//injected from build
	version = ""
	build   = ""
)

func main() {
	viper.Set("version", version)
	viper.Set("build", build)
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
