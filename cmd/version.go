package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitCommit is the git reference injected at build
	GitCommit string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	versionCmd    = &cobra.Command{
		Use:   "version",
		Short: "Show the version of chatbus",
		Run: RunWithCtx(func(ctx context.Context, cmd *cobra.Command, args []string) {
			fmt.Println("chatbus version: " + GitCommit)
		}),
	}
)
