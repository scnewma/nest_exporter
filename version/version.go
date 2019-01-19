package version

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime"
	"strings"
)

// Build information.
var (
	Program   = "nest_exporter"
	Version   string
	GoVersion = runtime.Version()
)

var versionInfoTmpl = `
{{ .program }}, version {{ .version }}
  go version:  {{ .goVersion }}
`

// Print returns basic version information suitable for
// display in a command line output.
func Print() string {
	m := map[string]string{
		"program":   Program,
		"version":   Version,
		"goVersion": GoVersion,
	}

	t := template.Must(template.New("version").Parse(versionInfoTmpl))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}

	return strings.TrimSpace(buf.String())
}

// Info returns shortened version information.
func Info() string {
	return fmt.Sprintf("(version=%s)", Version)
}
