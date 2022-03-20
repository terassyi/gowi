package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func insertMemData(mem []byte, offset int, data []byte) {
	for i, d := range data {
		mem[offset+i] = d
	}
}

func TestMemoryInitData(t *testing.T) {
	for _, d := range []struct {
		memory *Memory
		insert []struct {
			offset int32
			data   []byte
		}
	}{
		{
			memory: &Memory{Type: &types.MemoryType{Limits: &types.Limits{Min: 1}}, Data: make([]byte, PAGE_SIZE*1)},
			insert: []struct {
				offset int32
				data   []byte
			}{
				{offset: 0, data: []byte("a")},
			},
		},
		{
			memory: &Memory{Type: &types.MemoryType{Limits: &types.Limits{Min: 1}}, Data: make([]byte, PAGE_SIZE*1)},
			insert: []struct {
				offset int32
				data   []byte
			}{
				{offset: 0, data: []byte("a")},
				{offset: 100, data: []byte("abc")},
			},
		},
	} {
		exp := make([]byte, PAGE_SIZE*d.memory.Type.Limits.Min)
		for _, ins := range d.insert {
			insertMemData(exp, int(ins.offset), ins.data)
		}
		for _, ins := range d.insert {
			err := d.memory.initData(ins.offset, ins.data)
			require.NoError(t, err)
		}
		assert.Equal(t, exp, d.memory.Data)
	}
}

func TestMemoryInitData_Panic(t *testing.T) {
	for _, d := range []struct {
		memory *Memory
		insert []struct {
			offset int32
			data   []byte
			flag   bool
		}
	}{
		{
			memory: &Memory{Type: &types.MemoryType{Limits: &types.Limits{Min: 0}}, Data: make([]byte, PAGE_SIZE*0)},
			insert: []struct {
				offset int32
				data   []byte
				flag   bool
			}{
				{offset: 0, data: []byte("a"), flag: true},
			},
		},
		{
			memory: &Memory{Type: &types.MemoryType{Limits: &types.Limits{Min: 1, Max: 2}}, Data: make([]byte, PAGE_SIZE*1)},
			insert: []struct {
				offset int32
				data   []byte
				flag   bool
			}{
				{offset: 0, data: []byte("a"), flag: false},
				{offset: int32(PAGE_SIZE + 10), data: []byte("deadbeef"), flag: true},
			},
		},
	} {
		// exp := make([]byte, PAGE_SIZE*d.memory.Type.Limits.Min)
		for _, ins := range d.insert {
			err := d.memory.initData(ins.offset, ins.data)
			if ins.flag {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		}
	}
}
