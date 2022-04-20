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

func WithSize(v, f int) *Stack {
	if v == 0 && f == 0 {
		return New()
	}
	return &Stack{
		Value: &ValueStack{values: make([]value.Value, 0, v)},
		Frame: &FrameStack{frames: make([]Frame, 0, f)},
		Label: &LabelStack{labels: make([]Label, 0, f)},
	}
}

func WithValue(values []value.Value, frames []Frame, labels []Label) (*Stack, error) {
	stack := New()
	for _, v := range values {
		if err := stack.Value.push(v); err != nil {
			return nil, err
		}
	}
	for _, f := range frames {
		if err := stack.Frame.push(f); err != nil {
			return nil, err
		}
	}
	for _, l := range labels {
		if err := stack.Label.push(l); err != nil {
			return nil, err
		}
	}
	return stack, nil
}

func (s *Stack) Push(val value.Value) error {
	return s.Value.push(val)
}

func (s *Stack) Pop() (value.Value, error) {
	return s.Value.pop()
}

func (s *Stack) Top() (value.Value, error) {
	return s.Value.top()
}

func (s *Stack) PushValue(val value.Value) error {
	return s.Value.push(val)
}

func (s *Stack) PopValue() (value.Value, error) {
	poped := make([]value.Value, 0, 10)
	for s.Value.len() > 0 {
		val, err := s.Value.pop()
		if err != nil {
			return nil, fmt.Errorf("pop vaue: %w", err)
		}
		poped = append(poped, val)
		if val.ValType() == value.ValTypeNum {
			break
		}
	}
	val := poped[len(poped)-1]
	for i := len(poped) - 2; i > 0; i-- {
		if err := s.PushValue(poped[i]); err != nil {
			return nil, fmt.Errorf("pop value: %w", err)
		}
	}
	return val, nil
}

func (s *Stack) PopValues(n int) ([]value.Value, error) {
	values := make([]value.Value, 0, n)
	for i := 0; i < n; i++ {
		v, err := s.PopValue()
		if err != nil {
			return nil, fmt.Errorf("pop values: %w", err)
		}
		values = append(values, v)
	}
	return values, nil
}

func (s *Stack) PopValuesRev(n int) ([]value.Value, error) {
	values, err := s.PopValues(n)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(values)/2; i++ {
		values[i], values[len(values)-i-1] = values[len(values)-i-1], values[i]
	}
	return values, nil
}

func (s *Stack) TopValue() (value.Value, error) {
	val, err := s.Value.top()
	if err != nil {
		return nil, err
	}
	if val.ValType() != value.ValTypeNum {
		return s.TopValue()
	}
	return val, nil
}

func (s *Stack) Len() int {
	return s.Value.len()
}

func (s *Stack) RefValue() (value.Value, error) {
	if s.Value.isEmpty() {
		return nil, fmt.Errorf("ref: %w", StackIsEmpty)
	}
	for i := s.Value.len() - 1; i > 0; i-- {
		v := s.Value.values[i]
		if v.ValType() == value.ValTypeNum {
			return v, nil
		}
	}
	return nil, fmt.Errorf("ref: Value is not found")
}

func (s *Stack) RefNValue(n int) ([]value.Value, error) {
	if s.Value.len() < n {
		return nil, StackIsEmpty
	}
	values := make([]value.Value, 0, n)
	for i := s.Value.len() - 1; i > 0; i-- {
		v := s.Value.values[i]
		if v.ValType() == value.ValTypeNum {
			values = append(values, v)
		}
		if len(values) == n {
			break
		}
	}
	return values, nil
}

func (s *Stack) RefNValueRev(n int) ([]value.Value, error) {
	values, err := s.RefNValue(n)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(values)/2; i++ {
		values[i], values[len(values)-i-1] = values[len(values)-i-1], values[i]
	}
	return values, nil
}

func (s *Stack) ValidateValue(values []types.ValueType) error {
	for i, t := range values {
		val := s.Value.values[len(s.Value.values)-i-1]
		if val.ValType() == value.ValTypeNum {
			if !val.(value.Number).ValidateValueType(t) {
				return ValueStackTypeNotMatch
			}
		} else if val.ValType() == value.ValTypeVec {
			if t != types.V128 {
				return ValueStackTypeNotMatch
			}
		} else if val.ValType() == value.ValTypeFrame || val.ValType() == value.ValTypeLabel {
			continue
		} else {
			return ValueStackTypeNotMatch
		}
	}
	return nil
}

type ValueStack struct {
	values []value.Value
	sp     uint32
}

func (vs *ValueStack) push(val value.Value) error {
	if len(vs.values) >= VALUE_STACK_LIMIT {
		return fmt.Errorf("value stack push: %w", StackLimit)
	}
	vs.values = append(vs.values, val)
	return nil
}

func (vs *ValueStack) pop() (value.Value, error) {
	if len(vs.values) == 0 {
		return nil, fmt.Errorf("value stack pop: %w", StackIsEmpty)
	}
	val := vs.values[len(vs.values)-1]
	vs.values = vs.values[:len(vs.values)-1]
	return val, nil
}

func (vs *ValueStack) top() (value.Value, error) {
	if len(vs.values) == 0 {
		return nil, fmt.Errorf("value stack pop: %w", StackIsEmpty)
	}
	return vs.values[len(vs.values)-1], nil
}

func (vs *ValueStack) len() int {
	return len(vs.values)
}

func (vs *ValueStack) isEmpty() bool {
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

func (vs *ValueStack) popN(n int) ([]value.Value, error) {
	values := make([]value.Value, 0, n)
	for i := 0; i < n; i++ {
		v, err := vs.pop()
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}

func (vs *ValueStack) popNRev(n int) ([]value.Value, error) {
	values := make([]value.Value, n)
	for i := 0; i < n; i++ {
		v, err := vs.pop()
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

func (fs *FrameStack) push(frame Frame) error {
	if len(fs.frames) >= FRAME_STACK_LIMIT {
		return fmt.Errorf("frame stack push: %w", StackLimit)
	}
	fs.frames = append(fs.frames, frame)
	return nil
}

func (fs *FrameStack) pop() (*Frame, error) {
	if len(fs.frames) == 0 {
		return nil, fmt.Errorf("frame stack pop: %w", StackIsEmpty)
	}
	f := fs.frames[len(fs.frames)-1]
	fs.frames = fs.frames[:len(fs.frames)-1]
	return &f, nil
}

func (fs *FrameStack) top() (*Frame, error) {
	if len(fs.frames) == 0 {
		return nil, fmt.Errorf("frame stack top: %w", StackIsEmpty)
	}
	return &fs.frames[len(fs.frames)-1], nil
}

func (fs *FrameStack) len() int {
	return len(fs.frames)
}

func (fs *FrameStack) isEmpty() bool {
	return len(fs.frames) == 0
}

func (s *Stack) PushFrame(frame Frame) error {
	if err := s.PushValue(value.DummyFrame); err != nil {
		return fmt.Errorf("push frame: %w", err)
	}
	return s.Frame.push(frame)
}

func (s *Stack) PopFrame() (*Frame, error) {
	f, err := s.Frame.pop()
	if err != nil {
		return nil, fmt.Errorf("pop frame: %w", err)
	}
	frameIndex := 0
	for i := s.Value.len() - 1; i > 0; i-- {
		v := s.Value.values[i]
		if v.ValType() == value.ValTypeFrame {
			frameIndex = i
			break
		}
	}
	s.Value.values = append(s.Value.values[:frameIndex], s.Value.values[frameIndex+1:]...)
	return f, nil
}

func (s *Stack) TopFrame() (*Frame, error) {
	return s.Frame.top()
}

func (s *Stack) IsFrameEmpty() bool {
	return s.Frame.isEmpty()
}

func (s *Stack) LenFrame() int {
	return s.Frame.len()
}

type LabelStack struct {
	labels []Label
}

func (ls *LabelStack) push(label Label) error {
	if len(ls.labels) >= LABEL_STACK_LIMIT {
		return fmt.Errorf("label stack push: %w", StackLimit)
	}
	ls.labels = append(ls.labels, label)
	return nil
}

func (ls *LabelStack) pop() (*Label, error) {
	if len(ls.labels) == 0 {
		return nil, fmt.Errorf("label stack pop: %w", StackIsEmpty)
	}
	l := ls.labels[len(ls.labels)-1]
	ls.labels = ls.labels[:len(ls.labels)-1]
	return &l, nil
}

func (ls *LabelStack) top() (*Label, error) {
	if len(ls.labels) == 0 {
		return nil, fmt.Errorf("label stack Top: %w", StackIsEmpty)
	}
	return &ls.labels[len(ls.labels)-1], nil
}

func (ls *LabelStack) ref(n int) (*Label, error) {
	if ls.len() < n {
		return nil, fmt.Errorf("label ref: %w", InvalidStackLength)
	}
	return &ls.labels[len(ls.labels)-1-n], nil
}

func (ls *LabelStack) len() int {
	return len(ls.labels)
}

func (ls *LabelStack) IsEmpty() bool {
	return len(ls.labels) == 0
}

func (s *Stack) PushLabel(label Label) error {
	if err := s.PushValue(value.DummyLabel); err != nil {
		return fmt.Errorf("push label: %w", err)
	}
	return s.Label.push(label)
}

func (s *Stack) PopLabel() (*Label, error) {
	label, err := s.Label.pop()
	if err != nil {
		return nil, fmt.Errorf("label pop: %w", err)
	}
	var labelIndex int
	for i := s.Value.len() - 1; i > 0; i-- {
		v := s.Value.values[i]
		if v.ValType() == value.ValTypeLabel {
			labelIndex = i
			break
		}
	}
	s.Value.values = append(s.Value.values[:labelIndex], s.Value.values[labelIndex+1:]...)
	return label, nil
}

func (s *Stack) RefLabel(n int) (*Label, error) {
	return s.Label.ref(n)
}

func (s *Stack) TopLabel() (*Label, error) {
	return s.Label.top()
}

func (s *Stack) LenLabel() int {
	return s.Label.len()
}

func (s *Stack) IsLabelEmpty() bool {
	return s.Label.IsEmpty()
}

type Label struct {
	Instructions []instruction.Instruction
	N            uint8
	Sp           int
	Flag         bool
}

type Frame struct {
	Locals []value.Value
	Module *instance.Module
}
