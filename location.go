package justest

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"runtime"
	"slices"
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
	b, err := os.ReadFile(file)
	if err != nil {
		panic(fmt.Errorf("failed reading '%s': %w", file, err))
	}

	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", b, parser.ParseComments)
	if err != nil {
		panic(fmt.Errorf("failed parsing '%s': %w", file, err))
	}

	// Find the statement that contains the given line number
	var sourceCode bytes.Buffer
	ast.Inspect(f, func(n ast.Node) bool {
		// If the node is nil or the statement is not in the node, continue
		if n == nil || fileSet.Position(n.Pos()).Line > line || fileSet.Position(n.End()).Line < line {
			return true
		}

		// Check if the node is a function call expression - ignore otherwise
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Check that this is indeed the call to "Now()", "For(...)" or "Within(...)"
		if assertCE, ok := ce.Fun.(*ast.SelectorExpr); !ok {
			return true
		} else if !slices.Contains([]string{"Now", "For", "Within"}, assertCE.Sel.Name) {
			return true
		}

		if err := format.Node(&sourceCode, fileSet, n); err != nil {
			panic(err)
		}
		return false
	})

	// If not found in the AST, use the line from the file (next best thing...)
	if sourceCode.Len() == 0 {
		if lines := strings.Split(string(b), "\n"); len(lines) > line {
			sourceCode.WriteString(strings.TrimSpace(lines[line-1]))
		} else {
			return fmt.Sprintf("(missing) %s:%d", file, line)
		}
	}

	result := sourceCode.String()

	// Highlight the result if configured to
	if highlight {
		output := bytes.Buffer{}
		if err := quick.Highlight(&output, sourceCode.String(), "go", goSourceFormatter, goSourceStyle[displayMode]); err == nil {
			result = output.String()
		}
		result = output.String()
	}

	return result
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
