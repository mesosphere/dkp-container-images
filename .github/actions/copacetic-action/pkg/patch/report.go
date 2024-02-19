package patch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

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

func WriteMarkdown(report Report, w io.Writer) error {
	doc := md.NewMarkdown(w)

	imagesTable := md.TableSet{
		Header: []string{"Image", "Patched", "Error"},
	}

	details := [][]string{}

	for i, row := range report {
		mdRow := []string{
			md.Code(row.Image),
			md.Code(row.Patched),
		}
		if row.Error != "" {
			mdRow = append(mdRow, md.Link("View error", fmt.Sprintf("#error-%d", i)))
			details = append(details, []string{
				row.Image,
				row.Error,
				row.Output,
			})
		}
		imagesTable.Rows = append(imagesTable.Rows, mdRow)
	}

	doc.H2("Patched images").LF().Table(imagesTable)

	for _, detail := range details {
		doc.Details(detail[0], fmt.Sprintf("```%s```", detail[1]))
	}

	return doc.Build()
}
