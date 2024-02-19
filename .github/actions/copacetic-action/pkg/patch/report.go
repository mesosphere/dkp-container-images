package patch

import (
	"encoding/json"
	"errors"
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
	for _, row := range report {
		mdRow := []string{
			row.Image,
			row.Patched,
			row.Error,
		}
		for i := range mdRow {
			if len(mdRow[i]) > 0 {
				mdRow[i] = md.Code(mdRow[i])
			}
		}
		imagesTable.Rows = append(imagesTable.Rows, mdRow)
	}

	doc.H2("Patched images").PlainText("").Table(imagesTable)
	return doc.Build()
}
