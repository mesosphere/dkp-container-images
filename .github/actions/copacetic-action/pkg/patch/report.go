package patch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	md "github.com/go-spectest/markdown"

	"github.com/d2iq-labs/copacetic-action/pkg/image"
)

type Report []Item

type Item struct {
	Image   string `json:"image"`
	Patched string `json:"patched,omitempty"`
	Error   string `json:"error"`
	Output  string `json:"output,omitempty"`
}

func WriteJSON(tasks []*Task, w io.Writer) error {
	r := Report{}

	for _, task := range tasks {
		t := task
		item := Item{
			Image:   t.Image,
			Patched: t.Patch.Patched,
		}
		if task.Error != nil {
			item.Error = t.Error.Error()
			cmdErr := &image.CmdErr{}
			if errors.As(t.Error, &cmdErr) {
				item.Output = string(cmdErr.Output)
			}
		}

		r = append(r, item)
	}

	return json.NewEncoder(w).Encode(r)
}

func WriteMarkdown(ctx context.Context, report Report, w io.Writer, printCVEs bool) error {
	doc := md.NewMarkdown(w)

	imagesTable := md.TableSet{
		Header: []string{"Image", "Patched", "Error"},
	}

	errorDetails := [][]string{}

	for i, row := range report {
		mdRow := []string{
			md.Code(row.Image),
			md.Code(row.Patched),
		}

		if printCVEs {
			mdRow[0] = fmt.Sprintf("`%s`<br>%s", row.Image, scanImage(ctx, row.Image))
			mdRow[1] = fmt.Sprintf("`%s`<br>%s", row.Patched, scanImage(ctx, row.Patched))
		}

		if row.Error != "" {
			mdRow = append(mdRow, md.Link("View error", fmt.Sprintf("#error-%d", i)))

			detailsRow := []string{
				row.Image,
				fmt.Sprintf("error-%d", i),
				row.Error,
			}

			if len(row.Output) > 0 {
				detailsRow = append(detailsRow, row.Output)
			}

			errorDetails = append(errorDetails, detailsRow)
		} else {
			mdRow = append(mdRow, "")
		}

		imagesTable.Rows = append(imagesTable.Rows, mdRow)
	}

	doc.H2("Patched images").LF().Table(imagesTable)

	if len(errorDetails) > 0 {
		doc.H2("Errors")
	}
	for _, detail := range errorDetails {
		doc.PlainText(fmt.Sprintf(`<a name="%s"></a>`, detail[1]))
		detailsContent := []string{}
		for i := 2; i < len(detail); i++ {
			detailsContent = append(detailsContent, fmt.Sprintf("<pre>%s</pre>", detail[i]))
		}
		doc.Details(detail[0], strings.Join(detailsContent, "\n"))
	}

	return doc.Build()
}

func scanImage(ctx context.Context, imageRef string) string {
	report, err := image.Scan(ctx, imageRef, image.ScanAllOS)
	if err != nil {
		return md.Code(err.Error())
	}

	counts := map[string]int{
		"CRITICAL": 0,
		"HIGH":     0,
	}
	for _, vuln := range report.Vulnerabilities() {
		if _, ok := counts[vuln.Severity]; ok {
			counts[vuln.Severity] = counts[vuln.Severity] + 1
		}
	}

	parts := []string{}
	for severity, count := range counts {
		parts = append(parts, fmt.Sprintf("`%d` %s", count, md.Bold(severity)))
	}
	return strings.Join(parts, " ")
}
