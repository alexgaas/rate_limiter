package cmd

import (
	"errors"
	"fmt"
	"os"
	"quoter/internal"

	"github.com/spf13/cobra"
)

var (
	// some global variables
	configFile string
	debug      bool
	log        string

	g internal.CmdGlobal

	InitError error

	rootCmd = &cobra.Command{
		SilenceUsage: true,
		Use:          "quoter",
		Short:        "quoter service application",
		Long:         fmt.Sprintf(``),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if InitError != nil {
				os.Exit(1)
			}

			var Overrides internal.ConfigOverrides
			Overrides.Debug = debug
			Overrides.Log = log
			if g.Log != nil {
				g.Log.UpdateOverrides(Overrides)
			}
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(InitConfig)

	g = internal.CmdGlobal{Cmd: rootCmd, Opts: &internal.ConfYaml{}}

	// global command line switches, could override
	// configuration variables
	rootCmd.PersistentFlags().StringVarP(&configFile, "config file", "C", "",
		"configuration file")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d",
		false, "debug output, default: 'false'")
	rootCmd.PersistentFlags().StringVarP(&log, "log", "", "",
		"log file, default: 'stdout'")

	C, _ := internal.CreateCore(&g)

	// server sub-command
	serverCmd := cmdServer{g: &g}
	serverCmd.core = C
	rootCmd.AddCommand(serverCmd.Command())

	return
}

func InitConfig() {
	var err error
	var Overrides internal.ConfigOverrides
	Overrides.Debug = debug
	Overrides.Log = log

	*g.Opts, InitError = internal.LoadConf(configFile, Overrides)

	if InitError != nil {
		os.Exit(1)
		return
	}

	if g.Log, err = internal.CreateLog(g.Opts); err != nil {
		InitError = errors.New(fmt.Sprintf("logfile initialization error, err:'%s'", err))
		os.Exit(1)
		return
	}
}
