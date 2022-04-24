package value

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewI32_uint32(t *testing.T) {
	for _, d := range []struct {
		val uint32
		exp I32
	}{
		{val: 1, exp: I32(1)},
		{val: 0, exp: I32(0)},
		{val: 0xff, exp: I32(0xff)},
		{val: 0xffffffff, exp: I32(0xffffffff)},
		{val: 1 << 31, exp: I32(1 << 31)},
	} {
		require.Equal(t, d.exp, NewI32(d.val))
	}
}

func TestNewI32_int32(t *testing.T) {
	for _, d := range []struct {
		val int32
		exp I32
	}{
		{val: 0, exp: I32(0)},
		{val: 1, exp: I32(1)},
		{val: 0xff, exp: I32(0xff)},
		{val: 1 << 30, exp: I32(1 << 30)},
		{val: -1, exp: I32(0xffffffff)},
		{val: -3, exp: I32(0xfffffffd)},
		{val: -559038801, exp: I32(0xdeadbeaf)},
	} {
		require.Equal(t, d.exp, NewI32(d.val))
	}
}

func TestI32Signed(t *testing.T) {
	for _, d := range []struct {
		val I32
		exp int32
	}{
		{val: I32(0), exp: 0},
		{val: I32(1), exp: 1},
		{val: I32(0xef), exp: 0xef},
		{val: I32(0xeff1bfd5), exp: -269369387},
		{val: I32(0xffffffff), exp: -1},
	} {
		require.Equal(t, d.exp, d.val.Signed())
	}
}

func TestNew64_uint64(t *testing.T) {
	for _, d := range []struct {
		val uint64
		exp I64
	}{
		{val: 0, exp: I64(0)},
		{val: 1, exp: I64(1)},
		{val: 0xff, exp: I64(0xff)},
		{val: 1 << 30, exp: I64(1 << 30)},
		{val: 1 << 63, exp: I64(1 << 63)},
	} {
		require.Equal(t, d.exp, NewI64(d.val))
	}
}

func TestNew64_int64(t *testing.T) {
	for _, d := range []struct {
		val int64
		exp I64
	}{
		{val: 0, exp: I64(0)},
		{val: 1, exp: I64(1)},
		{val: 0xff, exp: I64(0xff)},
		{val: 1 << 30, exp: I64(1 << 30)},
		{val: 1 << 62, exp: I64(1 << 62)},
		{val: -1, exp: I64(0xffffffffffffffff)},
		{val: -3, exp: I64(0xfffffffffffffffd)},
		{val: -559038801, exp: I64(0xffffffffdeadbeaf)},
	} {
		require.Equal(t, d.exp, NewI64(d.val))
	}
}

func TestI64Signed(t *testing.T) {
	for _, d := range []struct {
		val I64
		exp int64
	}{
		{val: I64(0), exp: 0},
		{val: I64(1), exp: 1},
		{val: I64(1 << 30), exp: 1 << 30},
		{val: I64(0xffffffffffffffff), exp: -1},
		{val: I64(0xfffffffffffffffd), exp: -3},
	} {
		require.Equal(t, d.exp, d.val.Signed())
	}
}
