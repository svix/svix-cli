package validators

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NoArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("`%s` does not take any arguments. Run `%s --help` for usage information.",
				cmd.CommandPath(),
				cmd.CommandPath(),
			)
		}

		return nil
	}
}

func RangeArgs(min int, max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < min || len(args) > max {
			return fmt.Errorf("`%s` accepts between %d and %d arg(s), received %d. Run `%s --help` for usage information.",
				cmd.CommandPath(),
				min,
				max,
				len(args),
				cmd.CommandPath(),
			)
		}
		return nil
	}
}

func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n {
			msg := "`%s`requires %d args, received %d. Run `%s --help` for usage information."
			if n == 1 {
				msg = "`%s`requires 1 arg, received %d. Run `%s --help` for usage information."
			}
			return fmt.Errorf(msg,
				cmd.CommandPath(),
				n,
				len(args),
				cmd.CommandPath(),
			)
		}
		return nil
	}
}

func MinimumNArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < n {
			msg := "`%s` requires at least %d args, only received %d. Run `%s --help` for usage information."
			if n == 1 {
				msg = "`%s` requires at least 1 arg, only received %d. Run `%s --help` for usage information."
			}
			return fmt.Errorf(msg,
				cmd.CommandPath(),
				n,
				len(args),
				cmd.CommandPath(),
			)
		}
		return nil
	}
}
