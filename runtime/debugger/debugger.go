package debugger

type DebugLevel int

const (
	DebugLevelNoLog       DebugLevel = 0
	DebugLevelLogOnly     DebugLevel = 1
	DebugLevelShowContext DebugLevel = 2
	DebugLevelInterrupt   DebugLevel = 3
)

type Debugger struct {
	level DebugLevel
}

func New(level DebugLevel) (*Debugger, error) {
	return &Debugger{level: level}, nil
}
