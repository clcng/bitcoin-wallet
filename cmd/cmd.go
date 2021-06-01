package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/clcng/bitcoin-wallet/pkg/log"
	"github.com/clcng/bitcoin-wallet/pkg/server"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	configPath string

	RootCmd = &cobra.Command{
		Use:  "bitcoin-wallet",
		RunE: runHttpServ,
	}
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config", "config file name")
	RootCmd.PersistentFlags().StringVarP(&configPath, "config-path", "p", "./", "config file path")
	RootCmd.PersistentFlags().Bool("debug", true, "debug enabled")

	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

}

func initConfig() {
	viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		die("failure to load config file. %v", err)
	}
	initLogger()
	initService()
}

func initLogger() {
	if viper.GetBool("debug") {
		log.Level(zerolog.DebugLevel)
	} else {
		log.Level(zerolog.InfoLevel)
	}

	if viper.GetString("mode") == "dev" {
		log.Output(zerolog.ConsoleWriter{
			Out:     os.Stdout,
			NoColor: runtime.GOOS == "windows",
		})
	}
}

func initService() {
	log.Info().Msgf("[init service]version: %s, build: %s, mode: %s, debug %v",
		viper.GetString("version"),
		viper.GetString("build"),
		viper.GetString("mode"),
		viper.GetBool("debug"))
}

func runHttpServ(cmd *cobra.Command, args []string) error {
	go serv(cmd, args)

	//wait for exit signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}

func serv(cmd *cobra.Command, args []string) {
	addr := viper.GetString("addr")
	err := server.WalletServ(addr)
	if err != nil {
		die("failed to serve server %v", err)
	}
}
