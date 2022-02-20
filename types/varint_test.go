package types

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeUint32(t *testing.T) {
	for _, c := range []struct {
		input    VarUint32
		expected []byte
	}{
		{input: VarUint32(0), expected: []byte{0x00}},
		{input: VarUint32(1), expected: []byte{0x01}},
		{input: VarUint32(4), expected: []byte{0x04}},
		{input: VarUint32(16256), expected: []byte{0x80, 0x7f}},
		{input: VarUint32(624485), expected: []byte{0xe5, 0x8e, 0x26}},
		{input: VarUint32(165675008), expected: []byte{0x80, 0x80, 0x80, 0x4f}},
		{input: VarUint32(0xffffffff), expected: []byte{0xff, 0xff, 0xff, 0xff, 0xf}},
	} {
		require.Equal(t, c.expected, c.input.Encode())
	}
}

func TestDecodeUint32(t *testing.T) {
	for _, c := range []struct {
		bytes  []byte
		exp    VarUint32
		expErr bool
	}{
		{bytes: []byte{0xff, 0xff, 0xff, 0xff, 0xf}, exp: VarUint32(0xffffffff)},
		{bytes: []byte{0x00}, exp: VarUint32(0)},
		{bytes: []byte{0x04}, exp: VarUint32(4)},
		{bytes: []byte{0x01}, exp: VarUint32(1)},
		{bytes: []byte{0x80, 0x7f}, exp: VarUint32(16256)},
		{bytes: []byte{0xe5, 0x8e, 0x26}, exp: VarUint32(624485)},
		{bytes: []byte{0x80, 0x80, 0x80, 0x4f}, exp: VarUint32(165675008)},
		{bytes: []byte{0x83, 0x80, 0x80, 0x80, 0x80, 0x00}, expErr: true},
		{bytes: []byte{0x82, 0x80, 0x80, 0x80, 0x70}, expErr: true},
		{bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x00}, expErr: true},
	} {
		actual, num, err := DecodeVarUint32(bytes.NewReader(c.bytes))
		if c.expErr {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, c.exp, actual)
			assert.Equal(t, len(c.bytes), num)
		}
	}
}

func TestDecodeInt32(t *testing.T) {
	for i, c := range []struct {
		bytes  []byte
		exp    VarInt32
		expErr bool
	}{
		{bytes: []byte{0x13}, exp: VarInt32(19)},
		{bytes: []byte{0x00}, exp: VarInt32(0)},
		{bytes: []byte{0x04}, exp: VarInt32(4)},
		{bytes: []byte{0xFF, 0x00}, exp: VarInt32(127)},
		{bytes: []byte{0x81, 0x01}, exp: VarInt32(129)},
		{bytes: []byte{0x7f}, exp: VarInt32(-1)},
		{bytes: []byte{0x81, 0x7f}, exp: VarInt32(-127)},
		{bytes: []byte{0xFF, 0x7e}, exp: VarInt32(-129)},
		{bytes: []byte{0xff, 0xff, 0xff, 0xff, 0x0f}, expErr: true},
		{bytes: []byte{0xff, 0xff, 0xff, 0xff, 0x4f}, expErr: true},
		{bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x70}, expErr: true},
	} {
		actual, num, err := DecodeVarInt32(bytes.NewReader(c.bytes))
		if c.expErr {
			assert.Error(t, err, fmt.Sprintf("%d-th got value %d", i, actual))
		} else {
			assert.NoError(t, err, i)
			assert.Equal(t, c.exp, actual, i)
			assert.Equal(t, len(c.bytes), num, i)
		}
	}
}

func TestDecodeInt64(t *testing.T) {
	for _, c := range []struct {
		bytes []byte
		exp   VarInt64
	}{
		{bytes: []byte{0x00}, exp: VarInt64(0)},
		{bytes: []byte{0x04}, exp: VarInt64(4)},
		{bytes: []byte{0xFF, 0x00}, exp: VarInt64(127)},
		{bytes: []byte{0x81, 0x01}, exp: VarInt64(129)},
		{bytes: []byte{0x7f}, exp: VarInt64(-1)},
		{bytes: []byte{0x81, 0x7f}, exp: VarInt64(-127)},
		{bytes: []byte{0xFF, 0x7e}, exp: VarInt64(-129)},
		{bytes: []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x7f},
			exp: VarInt64(-9223372036854775808)},
	} {
		actual, num, err := DecodeVarInt64(bytes.NewReader(c.bytes))
		require.NoError(t, err)
		assert.Equal(t, c.exp, actual)
		assert.Equal(t, len(c.bytes), num)
	}
}
