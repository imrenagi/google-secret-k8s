package cmd

import (
	"fmt"

	vs "github.com/imrenagi/google-secret-k8s/agent-sidecar/version"
	"github.com/spf13/cobra"
)

// NewVersionCmd returns a new `version` command to be used as a sub-command to root
func NewVersionCmd() *cobra.Command {

	versionCmd := cobra.Command{
		Use:   "version",
		Short: fmt.Sprintf("Print version information"),
		Run: func(cmd *cobra.Command, args []string) {

			version := vs.GetVersion()
			fmt.Printf("%s: %s\n", cliName, version)
			fmt.Printf("  BuildDate: %s\n", version.BuildDate)
			fmt.Printf("  GitCommit: %s\n", version.GitCommit)
			fmt.Printf("  GitTreeState: %s\n", version.GitTreeState)
			if version.GitTag != "" {
				fmt.Printf("  GitTag: %s\n", version.GitTag)
			}
			fmt.Printf("  GoVersion: %s\n", version.GoVersion)
			fmt.Printf("  Compiler: %s\n", version.Compiler)
			fmt.Printf("  Platform: %s\n", version.Platform)
		},
	}

	return &versionCmd
}
