package cmd

import (
	"fmt"
	"quoter/internal"

	"github.com/spf13/cobra"
)

type cmdServer struct {
	g *internal.CmdGlobal

	debug bool

	// back reference for core created (filled in root command)
	core *internal.Core
}

func (c *cmdServer) Command() *cobra.Command {

	// Main subcommand:server
	cmd := &cobra.Command{}
	cmd.Use = "server --- "
	cmd.Short = "Controlling server process"
	cmd.Long = `Description:
  ratelimiter service is started as a process under systemd control or as
  stand-alone server process.
`
	cmd.PersistentFlags().BoolVarP(&c.debug, "debug", "d", false, "debug output")

	// starting ratelimiter server: "ratelimiter server start"
	serverStartCmd := cmdServerStart{g: c.g}
	serverStartCmd.s = c
	cmd.AddCommand(serverStartCmd.Command())

	// stopping ratelimiter server: "ratelimiter server stop"
	serverStopCmd := cmdServerStop{g: c.g}
	serverStopCmd.s = c
	cmd.AddCommand(serverStopCmd.Command())

	return cmd
}

func (c *cmdServer) InitCommand(cmd *cobra.Command, args []string) (string, error) {
	// Getting item name (if supplied)
	item := ""
	if len(args) > 0 {
		item = args[0]
	}

	return item, nil
}

// command: "ratelimiter server start"
type cmdServerStart struct {
	g *internal.CmdGlobal
	s *cmdServer

	detach bool
}

func (c *cmdServerStart) Command() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "start"
	cmd.Short = "Starting server"
	cmd.Long = "Starting server"

	// detach should be used in non-systemd run
	cmd.PersistentFlags().BoolVarP(&c.detach, "detach", "D", false, "detach process")

	cmd.RunE = c.Run
	return cmd
}

func (c *cmdServerStart) Run(cmd *cobra.Command, args []string) error {
	C, err := internal.CreateCore(c.g)
	if err != nil {
		c.g.Log.Error(fmt.Sprintf("error creating core, err:'%s'", err))
		return err
	}

	var Overrides internal.CoreOverrides
	Overrides.Detach = c.detach
	if err := C.StartCore(&Overrides); err != nil {
		c.g.Log.Error(fmt.Sprintf("error starting core, err:'%s'", err))
		return err
	}

	return nil
}

// command: "ratelimiter server stop"
type cmdServerStop struct {
	g *internal.CmdGlobal
	s *cmdServer

	detach bool
}

func (c *cmdServerStop) Command() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "stop"
	cmd.Short = "Stopping server"
	cmd.Long = "Stopping server"

	// detach should be used in non-systemd run
	cmd.PersistentFlags().BoolVarP(&c.detach, "detach", "D", false, "detach process")

	cmd.RunE = c.Run
	return cmd
}

func (c *cmdServerStop) Run(cmd *cobra.Command, args []string) error {
	C, err := internal.CreateCore(c.g)
	if err != nil {
		c.g.Log.Error(fmt.Sprintf("error creating core, err:'%s'", err))
		return err
	}

	var Overrides internal.CoreOverrides
	Overrides.Detach = c.detach
	if err := C.StopCore(&Overrides); err != nil {
		c.g.Log.Error(fmt.Sprintf("error stopping core, err:'%s'", err))
		return err
	}

	return nil
}
