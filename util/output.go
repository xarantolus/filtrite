package util

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	releaseFile = "release.md"

	lineTemplate = "* {{if .Repo}}[`{{.ListName}}`](https://github.com/{{.Repo}}/releases/latest/download/{{.ListName}}.dat){{else}}`{{.ListName}}`{{end}}: updated {{.GotCount}}/{{.FullCount}} list{{if ne .FullCount 1}}s{{end}}{{if ne .ErrorCount 0}}, {{if eq .ErrorCount 1}}one error{{else}}{{.ErrorCount}} errors{{end}}{{end}}\n"
)

var (
	tmpl = template.Must(template.New("").Parse(lineTemplate))
)

type info struct {
	ListName string

	Repo string

	GotCount  int
	FullCount int

	ErrorCount int
}

// AppendReleaseList appends information about this list/download to the release file
func AppendReleaseList(fn string, gotCount, fullCount int) (err error) {
	f, err := os.OpenFile(releaseFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	var friendlyName = strings.TrimSuffix(filepath.Base(fn), ".txt")

	return tmpl.Execute(f, info{
		ListName:   friendlyName,
		Repo:       os.Getenv("GITHUB_REPOSITORY"),
		GotCount:   gotCount,
		FullCount:  fullCount,
		ErrorCount: fullCount - gotCount,
	})
}
