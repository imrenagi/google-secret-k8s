package cmd

import (
	"github.com/spf13/cobra"
)

var (
	cliName = "gsecret-agent-injector"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

// NewRootCommand register all command group
func NewRootCommand(args []string) *cobra.Command {

	var command = &cobra.Command{
		Use:   cliName,
		Short: "gsecret-agent-injector is injector for google secret agent",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}

	flags := command.PersistentFlags()

	command.AddCommand(
		NewVersionCmd(),
		NewServerCmd(),
	)

	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.Parse(args)

	return command
}
