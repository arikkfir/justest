package justest

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/arikkfir/justest/internal"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Display mode
type displayModeType string

const (
	displayModeLight displayModeType = "light"
	displayModeDark  displayModeType = "dark"
)

var (
	displayMode = displayModeLight
)

// Source code highlighting
const (
	goSourceFormatter string = "terminal256"
)

var (
	goSourceStyle = map[displayModeType]string{
		displayModeLight: "autumn",
		displayModeDark:  "catppuccin-mocha",
	}
	ignoredStackTracePrefixes = []string{
		"testing.",
		"github.com/arikkfir/justest/",
		"github.com/arikkfir/justest.",
	}
)

type Location struct {
	Function string
	File     string
	Line     int
	Source   string
}

//go:noinline
func nearestLocation() Location {
	l := Location{
		Function: "unknown",
		File:     "unknown",
		Line:     0,
		Source:   "<could not read source>",
	}
	for _, frame := range internal.CallStackAt(0) {
		function, file, line := frame.Location()

		startsWithAnIgnoredPrefix := false
		for _, prefix := range ignoredStackTracePrefixes {
			if strings.HasPrefix(function, prefix) {
				startsWithAnIgnoredPrefix = true
				break
			}
		}

		if !startsWithAnIgnoredPrefix {
			l.Function, l.File = function, file
			l.Line = line
			l.Source = readSourceAt(l.File, l.Line)
			break
		}
	}
	return l
}

//go:noinline
func readSourceAt(file string, line int) string {
	source := "<could not read source>"
	if b, err := os.ReadFile(file); err == nil {
		fileContents := string(b)
		lines := strings.Split(fileContents, "\n")
		if len(lines) > line {
			source = strings.TrimSpace(lines[line-1])
			output := bytes.Buffer{}
			if err := quick.Highlight(&output, source, "go", goSourceFormatter, goSourceStyle[displayMode]); err == nil {
				source = output.String()
			}
		}
	}
	return source
}

const (
	appleScriptDarkModeQuery string = `tell application "System Events" to tell appearance preferences to get dark mode`
)

func init() {
	cmd := exec.Command("osascript", "-e", appleScriptDarkModeQuery)
	if out, err := cmd.Output(); err != nil {
		displayMode = displayModeLight
		fmt.Printf("Error determining system's dark mode: %+v\n", err)
	} else if dark, err := strconv.ParseBool(strings.TrimSpace(string(out))); err != nil {
		displayMode = displayModeLight
		fmt.Printf("Error determining system's dark mode: %+v\n", err)
	} else if dark {
		displayMode = displayModeDark
	} else {
		displayMode = displayModeLight
	}
}
