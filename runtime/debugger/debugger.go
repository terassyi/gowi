package debugger

import (
	"fmt"
	"io"
	"os"

	"github.com/terassyi/gowi/runtime/value"
)

type DebugLevel int

const (
	DebugLevelNoLog         DebugLevel = 0
	DebugLevelLogOnly       DebugLevel = 1
	DebugLevelLogOnlyStdout DebugLevel = 2
	DebugLevelShowContext   DebugLevel = 3
	DebugLevelInterrupt     DebugLevel = 4
)

type Debugger struct {
	level  DebugLevel
	writer io.Writer
}

func New(level DebugLevel) *Debugger {
	switch level {
	case DebugLevelLogOnly:
		return &Debugger{level: level, writer: os.Stderr}
	case DebugLevelLogOnlyStdout:
		return &Debugger{level: level, writer: os.Stdout}
	default:
		return &Debugger{level: level, writer: io.Discard}
	}
}

func (d *Debugger) ShowResult(results []value.Value) {
	fmt.Fprintf(d.writer, "Execution Result = (")
	for i, res := range results {
		fmt.Fprintf(d.writer, "%v", res)
		if i < len(results)-1 {
			fmt.Fprintf(d.writer, ",")
		}
	}
	fmt.Fprintf(d.writer, ")")
}
