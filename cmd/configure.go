package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	svix "github.com/svixhq/svix-libs/go"
)

type configureCmd struct {
	cmd *cobra.Command
	sc  *svix.Svix
}

func newConfigureCmd(s *svix.Svix) *configureCmd {
	cc := &configureCmd{
		sc: s,
	}
	cc.cmd = &cobra.Command{
		Use:   "configure",
		Short: "Interactively configure you're Svix API credentials",
		Args:  cobra.ExactArgs(0),
		Run:   cc.run,
	}
	return cc
}

func (c *configureCmd) run(cmd *cobra.Command, args []string) {
	fmt.Println("configure")
}

func (c *configureCmd) Cmd() *cobra.Command {
	return c.cmd
}
