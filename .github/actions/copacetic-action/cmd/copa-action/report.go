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

var printCVEs = false

func NewMarkdownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "markdown PATH | -",
		Short: "Generate markdown report from JSON output",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			w := io.MultiWriter(os.Stdout, os.Stderr)
			return patch.WriteMarkdown(cmd.Context(), report, w, printCVEs)
		},
	}
	cmd.Flags().BoolVar(&printCVEs, "print-cves", printCVEs, "enable scanning and printing number of Critical and High CVEs")
	return cmd
}
