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
			return fmt.Errorf("`%s` accepts %d arg(s), received %d. Run `%s --help` for usage information.",
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
			return fmt.Errorf("`%s` requires at least %d arg(s), only received %d. Run `%s --help` for usage information.",
				cmd.CommandPath(),
				n,
				len(args),
				cmd.CommandPath(),
			)
		}
		return nil
	}
}
