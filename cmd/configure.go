package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/cfg"
	svix "github.com/svixhq/svix-libs/go"
)

type configureCmd struct {
	cmd *cobra.Command
	sc  *svix.Svix
	cfg *cfg.Config
}

func newConfigureCmd(cfg *cfg.Config, s *svix.Svix) *configureCmd {
	cc := &configureCmd{
		sc:  s,
		cfg: cfg,
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
