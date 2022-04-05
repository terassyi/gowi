package stack

import (
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

const (
	VALUE_STACK_LIMIT = 1024 * 1024
	FRAME_STACK_LIMIT = 64 * 1024
	LABEL_STACK_LIMIT = 64 * 1024
)

var (
	StackLimit             error = errors.New("Stack limit")
	StackIsEmpty           error = errors.New("Stack is empty")
	InvalidStackLength     error = errors.New("Invalid stack length")
	ValueStackTypeNotMatch error = errors.New("Value type in the stack is not matched")
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

func WithValue(values []value.Value, frames []Frame, labels []Label) (*Stack, error) {
	stack := New()
	for _, v := range values {
		if err := stack.Value.Push(v); err != nil {
			return nil, err
		}
	}
	for _, f := range frames {
		if err := stack.Frame.Push(f); err != nil {
			return nil, err
		}
	}
	for _, l := range labels {
		if err := stack.Label.Push(l); err != nil {
			return nil, err
		}
	}
	return stack, nil
}

func (s *Stack) PushValue(val value.Value) error {
	if err := s.Value.Push(val); err != nil {
		return err
	}
	label, err := s.Label.Top()
	if err != nil {
		return err
	}
	label.ValCounter++
	return nil
}

func (s *Stack) PopValue() (value.Value, error) {
	val, err := s.Value.Pop()
	if err != nil {
		return nil, err
	}
	label, err := s.Label.Top()
	if err != nil {
		return nil, err
	}
	label.ValCounter--
	return val, nil
}

func (s *Stack) PopValues(n int) ([]value.Value, error) {
	values, err := s.Value.PopN(n)
	if err != nil {
		return nil, err
	}
	label, err := s.Label.Top()
	if err != nil {
		return nil, err
	}
	label.ValCounter -= uint(n)
	return values, nil
}

func (s *Stack) PopValuesRev(n int) ([]value.Value, error) {
	values, err := s.Value.PopNRev(n)
	if err != nil {
		return nil, err
	}
	label, err := s.Label.Top()
	if err != nil {
		return nil, err
	}
	label.ValCounter -= uint(n)
	return values, nil
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
	val := vs.values[len(vs.values)-1]
	vs.values = vs.values[:len(vs.values)-1]
	return val, nil
}

func (vs *ValueStack) Top() (value.Value, error) {
	if len(vs.values) == 0 {
		return nil, fmt.Errorf("value stack pop: %w", StackIsEmpty)
	}
	return vs.values[len(vs.values)-1], nil
}

func (vs *ValueStack) Len() int {
	return len(vs.values)
}

func (vs *ValueStack) IsEmpty() bool {
	return len(vs.values) == 0
}

func (vs *ValueStack) Validate(ts []types.ValueType) error {
	for i, t := range ts {
		val := vs.values[len(vs.values)-i-1]
		if val.ValType() == value.ValTypeNum {
			if !val.(value.Number).ValidateValueType(t) {
				return ValueStackTypeNotMatch
			}
		} else if val.ValType() == value.ValTypeVec {
			if t != types.V128 {
				return ValueStackTypeNotMatch
			}
		} else {
			return ValueStackTypeNotMatch
		}
	}
	return nil
}

func (vs *ValueStack) PopN(n int) ([]value.Value, error) {
	values := make([]value.Value, 0, n)
	for i := 0; i < n; i++ {
		v, err := vs.Pop()
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func (vs *ValueStack) PopNRev(n int) ([]value.Value, error) {
	values := make([]value.Value, n)
	for i := 0; i < n; i++ {
		v, err := vs.Pop()
		if err != nil {
			return nil, err
		}
		values[n-i-1] = v
	}
	return values, nil
}

func (vs *ValueStack) RefNRev(n int) ([]value.Value, error) {
	if len(vs.values) < n {
		return nil, StackIsEmpty
	}
	values := make([]value.Value, n)
	for i := 0; i < n; i++ {
		v := vs.values[len(vs.values)-i-1]
		values[n-i-1] = v
	}
	return values, nil
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

func (fs *FrameStack) Top() (*Frame, error) {
	if len(fs.frames) == 0 {
		return nil, fmt.Errorf("frame stack top: %w", StackIsEmpty)
	}
	return &fs.frames[len(fs.frames)-1], nil
}

func (fs *FrameStack) Len() int {
	return len(fs.frames)
}

func (fs *FrameStack) IsEmpty() bool {
	return len(fs.frames) == 0
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

func (ls *LabelStack) Top() (*Label, error) {
	if len(ls.labels) == 0 {
		return nil, fmt.Errorf("label stack Top: %w", StackIsEmpty)
	}
	return &ls.labels[len(ls.labels)-1], nil
}

func (ls *LabelStack) Ref(n int) (*Label, error) {
	if ls.Len() < n {
		return nil, fmt.Errorf("label ref: %w", InvalidStackLength)
	}
	return &ls.labels[len(ls.labels)-1-n], nil
}

func (ls *LabelStack) Len() int {
	return len(ls.labels)
}

func (ls *LabelStack) IsEmpty() bool {
	return len(ls.labels) == 0
}

type Label struct {
	Instructions []instruction.Instruction
	N            uint8
	Sp           int
	Flag         bool
	ValCounter   uint
}

type Frame struct {
	Locals []value.Value
	Module *instance.Module
}
