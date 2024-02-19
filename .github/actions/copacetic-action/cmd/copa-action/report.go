package copa_action

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/d2iq-labs/copacetic-action/pkg/cli"
	"github.com/d2iq-labs/copacetic-action/pkg/patch"
	"github.com/spf13/cobra"
)

func NewMarkdownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markdown PATH | -",
		Short: "Generate markdown report from JSON output",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			input, err := cli.OpenFileOrStdin(args[0])
			if err != nil {
				return err
			}
			data, err := io.ReadAll(input)
			if err != nil {
				return err
			}

			report := patch.Report{}
			if err := json.Unmarshal(data, &report); err != nil {
				return fmt.Errorf("failed to read JSON report: %w", err)
			}

			return patch.WriteMarkdown(report, os.Stdout)
		},
	}
	return cmd
}
