package instance

import "github.com/terassyi/gowi/runtime/value"

type Export struct {
	Name  string
	Value value.Value // external val
}
