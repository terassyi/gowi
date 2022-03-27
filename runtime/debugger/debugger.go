package debugger

import (
	"fmt"
	"io"
	"os"

	"github.com/terassyi/gowi/runtime/value"
)

type DebugLevel int

const (
	DebugLevelNoLog       DebugLevel = 0
	DebugLevelLogOnly     DebugLevel = 1
	DebugLevelShowContext DebugLevel = 2
	DebugLevelInterrupt   DebugLevel = 3
)

type Debugger struct {
	level  DebugLevel
	writer io.Writer
}

func New(level DebugLevel) *Debugger {
	return &Debugger{level: level, writer: os.Stdout}
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
