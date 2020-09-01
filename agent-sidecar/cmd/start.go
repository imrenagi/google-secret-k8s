package cmd

import (
	"github.com/spf13/cobra"
)

// NewStartCmd returns start agent command
func NewStartCmd() *cobra.Command {

	startCmd := cobra.Command{
		Use:   "start",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	// startCmd.Flags().StringVar(&agentImage, "agent-image", "imrenagi/gsecret-agent:latest", "Agent image used as init or sidecar container")
	// startCmd.Flags().BoolVar(&requireAnnotation, "require-annotation", true, "If it is true, annotation should be given so that sidecar can be injected")

	return &startCmd
}
