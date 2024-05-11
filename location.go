package justest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"

	"github.com/arikkfir/justest/internal"
)

// Display mode
type displayModeType string

const (
	displayModeLight displayModeType = "light"
	displayModeDark  displayModeType = "dark"
)

var (
	displayMode = displayModeLight
	highlight   = true
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
			if highlight {
				output := bytes.Buffer{}
				if err := quick.Highlight(&output, source, "go", goSourceFormatter, goSourceStyle[displayMode]); err == nil {
					source = output.String()
				}
			}
		}
	}
	return source
}

func init() {
	const appleScriptDarkModeQuery string = `tell application "System Events" to tell appearance preferences to get dark mode`

	if highlightEnv := os.Getenv("JUSTEST_DISABLE_SOURCE_HIGHLIGHT"); highlightEnv != "" {
		if val, err := strconv.ParseBool(highlightEnv); err != nil {
			panic(fmt.Sprintf("Error parsing JUSTEST_HIGHLIGHT_SOURCE environment variable - illegal value: %s", highlightEnv))
		} else {
			highlight = val
		}
	}

	if highlight {
		switch runtime.GOOS {
		case "darwin":
			cmd := exec.Command("osascript", "-e", appleScriptDarkModeQuery)
			if out, err := cmd.Output(); err != nil {
				fmt.Printf("Error determining system's dark mode: %+v\n", err)
				highlight = false
			} else if dark, err := strconv.ParseBool(strings.TrimSpace(string(out))); err != nil {
				fmt.Printf("Error determining system's dark mode: %+v\n", err)
				highlight = false
			} else if dark {
				displayMode = displayModeDark
			} else {
				displayMode = displayModeLight
			}
		case "windows", "linux":
			// TODO: Implement a similar mechanism for Windows and Linux
			highlight = false
		}
	}
}
