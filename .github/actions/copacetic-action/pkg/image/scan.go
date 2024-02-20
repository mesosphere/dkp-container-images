package image

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/aquasecurity/trivy/pkg/types"
)

type Report struct {
	types.Report
}

func (r *Report) WriteTo(path string) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o755)
}

func (r *Report) Vulnerabilities() []types.DetectedVulnerability {
	vulnerabilities := []types.DetectedVulnerability{}
	for _, resultClass := range r.Results {
		vulnerabilities = append(vulnerabilities, resultClass.Vulnerabilities...)
	}
	return vulnerabilities
}

func VulnerabilitiesIdsSorted(vulns []types.DetectedVulnerability) []string {
	result := []string{}
	for _, v := range vulns {
		id := fmt.Sprintf("%s-%s-%s", v.VulnerabilityID, v.PkgName, v.InstalledVersion)
		result = append(result, id)
	}
	slices.Sort(result)
	return result
}

type CmdErr struct {
	Err    error
	Stdout []byte
	Stderr []byte
	Output []byte
}

func (e *CmdErr) Error() string {
	return e.Err.Error()
}

var (
	ScanFixableOS = []string{"--vuln-type", "os", "--ignore-unfixed"}
	ScanAllOS     = []string{"--vuln-type", "os"}
)

// Scan runs a trivy scan of a image and returns back report.
func Scan(ctx context.Context, imageName string, scanType []string) (*Report, error) {
	flags := append([]string{"image"}, scanType...)
	flags = append(flags, "--format", "json", imageName)
	cmd, stdout, stderr := prepareCmd(
		ctx, "trivy", flags...,
	)
	err := cmd.Run()
	if err != nil {
		return nil, &CmdErr{
			Err:    err,
			Stdout: stdout.Bytes(),
			Stderr: stderr.Bytes(),
		}
	}

	report := &Report{}
	if err := json.Unmarshal(stdout.Bytes(), report); err != nil {
		return nil, err
	}

	return report, nil
}

func prepareCmd(ctx context.Context, name string, args ...string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	cmd := exec.CommandContext(ctx, name, args...)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd, stdout, stderr
}
