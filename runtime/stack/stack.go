package stack

import (
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/value"
)

const (
	VALUE_STACK_LIMIT = 1024 * 1024
	FRAME_STACK_LIMIT = 64 * 1024
	LABEL_STACK_LIMIT = 64 * 1024
)

var (
	StackLimit   error = errors.New("Stack limit")
	StackIsEmpty error = errors.New("Stack is empty")
)

// https://webassembly.github.io/spec/core/exec/runtime.html#stack
type Stack struct {
	Value *ValueStack
	Frame *FrameStack
	Label *LabelStack
}

func New() *Stack {
	return &Stack{
		Value: &ValueStack{values: make([]value.Value, 0, VALUE_STACK_LIMIT)},
		Frame: &FrameStack{frames: make([]Frame, 0, FRAME_STACK_LIMIT)},
		Label: &LabelStack{labels: make([]Label, 0, LABEL_STACK_LIMIT)},
	}
}

type ValueStack struct {
	values []value.Value
	sp     uint32
}

func (vs *ValueStack) Push(val value.Value) error {
	if len(vs.values) >= VALUE_STACK_LIMIT {
		return fmt.Errorf("value stack push: %w", StackLimit)
	}
	vs.values = append(vs.values, val)
	return nil
}

func (vs *ValueStack) Pop() (value.Value, error) {
	if len(vs.values) == 0 {
		return nil, fmt.Errorf("value stack pop: %w", StackIsEmpty)
	}
	fmt.Println(len(vs.values))
	val := vs.values[len(vs.values)-1]
	vs.values = vs.values[:len(vs.values)-1]
	fmt.Println(len(vs.values))
	return val, nil
}

type FrameStack struct {
	frames []Frame
}

func (fs *FrameStack) Push(frame Frame) error {
	if len(fs.frames) >= FRAME_STACK_LIMIT {
		return fmt.Errorf("frame stack push: %w", StackLimit)
	}
	fs.frames = append(fs.frames, frame)
	return nil
}

func (fs *FrameStack) Pop() (*Frame, error) {
	if len(fs.frames) == 0 {
		return nil, fmt.Errorf("frame stack pop: %w", StackIsEmpty)
	}
	f := fs.frames[len(fs.frames)-1]
	fs.frames = fs.frames[:len(fs.frames)-1]
	return &f, nil
}

type LabelStack struct {
	labels []Label
}

func (ls *LabelStack) Push(label Label) error {
	if len(ls.labels) >= LABEL_STACK_LIMIT {
		return fmt.Errorf("label stack push: %w", StackLimit)
	}
	ls.labels = append(ls.labels, label)
	return nil
}

func (ls *LabelStack) Pop() (*Label, error) {
	if len(ls.labels) == 0 {
		return nil, fmt.Errorf("label stack pop: %w", StackIsEmpty)
	}
	l := ls.labels[len(ls.labels)-1]
	ls.labels = ls.labels[:len(ls.labels)-1]
	return &l, nil
}

type Label struct {
	Instructions []instruction.Instruction
	n            uint8
}

type Frame struct {
	Locals []value.Value
	Module *instance.Module
}
