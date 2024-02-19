package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	copa_action "github.com/d2iq-labs/copacetic-action/cmd/copa-action"
)

var rootCmd = &cobra.Command{
	Use:   "copa-action",
	Short: "Copacetic Action",
}

func Execute() {
	rootCmd.AddCommand(copa_action.NewPatchCmd())
	rootCmd.AddCommand(copa_action.NewMarkdownCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
