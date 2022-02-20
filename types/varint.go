package types

import (
	"errors"
	"fmt"
	"io"
)

var encodeCache = [0x80][]byte{
	{0x00}, {0x01}, {0x02}, {0x03}, {0x04}, {0x05}, {0x06}, {0x07}, {0x08}, {0x09}, {0x0a}, {0x0b}, {0x0c}, {0x0d}, {0x0e}, {0x0f},
	{0x10}, {0x11}, {0x12}, {0x13}, {0x14}, {0x15}, {0x16}, {0x17}, {0x18}, {0x19}, {0x1a}, {0x1b}, {0x1c}, {0x1d}, {0x1e}, {0x1f},
	{0x20}, {0x21}, {0x22}, {0x23}, {0x24}, {0x25}, {0x26}, {0x27}, {0x28}, {0x29}, {0x2a}, {0x2b}, {0x2c}, {0x2d}, {0x2e}, {0x2f},
	{0x30}, {0x31}, {0x32}, {0x33}, {0x34}, {0x35}, {0x36}, {0x37}, {0x38}, {0x39}, {0x3a}, {0x3b}, {0x3c}, {0x3d}, {0x3e}, {0x3f},
	{0x40}, {0x41}, {0x42}, {0x43}, {0x44}, {0x45}, {0x46}, {0x47}, {0x48}, {0x49}, {0x4a}, {0x4b}, {0x4c}, {0x4d}, {0x4e}, {0x4f},
	{0x50}, {0x51}, {0x52}, {0x53}, {0x54}, {0x55}, {0x56}, {0x57}, {0x58}, {0x59}, {0x5a}, {0x5b}, {0x5c}, {0x5d}, {0x5e}, {0x5f},
	{0x60}, {0x61}, {0x62}, {0x63}, {0x64}, {0x65}, {0x66}, {0x67}, {0x68}, {0x69}, {0x6a}, {0x6b}, {0x6c}, {0x6d}, {0x6e}, {0x6f},
	{0x70}, {0x71}, {0x72}, {0x73}, {0x74}, {0x75}, {0x76}, {0x77}, {0x78}, {0x79}, {0x7a}, {0x7b}, {0x7c}, {0x7d}, {0x7e}, {0x7f},
}

var (
	Overflow32Error error = errors.New("overflow 32 bit integer.")
	Overflow64Error error = errors.New("overflow 64 bit integer.")
)

const (
	MAX_VARINT32_LENGTH = 5
	MAX_VARINT64_LENGTH = 10
)

type VarUint1 uint8

type VarUint7 uint8

type VarUint32 uint32

func (u VarUint32) Encode() (buf []byte) {
	n := u
	if n < 0x80 {
		return encodeCache[n]
	}
	for {
		b := uint8(n & 0x7f)
		n = n >> 7
		if n != 0 {
			b |= 0x80
		}
		buf = append(buf, b)
		if b&0x80 == 0 {
			return buf
		}
	}
}

func DecodeVarUint32(reader io.Reader) (VarUint32, int, error) {
	var s uint32
	var ret uint32 = 0
	for i := 0; i < MAX_VARINT32_LENGTH; i++ {
		b, err := readByte(reader)
		if err != nil {
			return 0, 0, err
		}
		if b < 0x80 {
			// Unused bits must be all zero.
			if i == MAX_VARINT32_LENGTH-1 && (b&0xf0) > 0 {
				return 0, 0, Overflow32Error
			}
			return VarUint32(ret | uint32(b)<<s), i + 1, nil
		}
		ret |= (uint32(b) & 0x7f) << s
		s += 7
	}
	return 0, 0, Overflow32Error

}

type VarInt7 int8

type VarInt32 int32

func DecodeVarInt32(r io.Reader) (VarInt32, int, error) {
	var shift int
	var ret int32 = 0
	read := 0
	for {
		b, err := readByte(r)
		if err != nil {
			return 0, 0, fmt.Errorf("readByte failed: %w", err)
		}
		ret |= (int32(b) & 0x7f) << shift
		shift += 7
		read++
		if b&0x80 == 0 {
			if shift < 32 && (b&0x40) != 0 {
				ret |= ^0 << shift
			}
			// Over flow checks.
			// fixme: can be optimized.
			if read > 5 {
				return 0, 0, Overflow32Error
			} else if unused := b & 0b00110000; read == 5 && ret < 0 && unused != 0b00110000 {
				return 0, 0, Overflow32Error
			} else if read == 5 && ret >= 0 && unused != 0x00 {
				return 0, 0, Overflow32Error
			}
			return VarInt32(ret), read, nil
		}
	}
}

type VarInt64 int64

func DecodeVarInt64(r io.Reader) (VarInt64, int, error) {
	const (
		int64Mask3 = 1 << 6
		int64Mask4 = ^0
	)
	var shift int
	var ret int64 = 0
	read := 0
	for {
		b, err := readByte(r)
		if err != nil {
			return 0, 0, fmt.Errorf("readByte failed: %w", err)
		}
		ret |= (int64(b) & 0x7f) << shift
		shift += 7
		read++
		if b&0x80 == 0 {
			if shift < 64 && (b&int64Mask3) == int64Mask3 {
				ret |= int64Mask4 << shift
			}
			// Over flow checks.
			// fixme: can be optimized.
			if read > 10 {
				return 0, 0, Overflow64Error
			} else if unused := b & 0b00111110; read == 10 && ret < 0 && unused != 0b00111110 {
				return 0, 0, Overflow64Error
			} else if read == 10 && ret >= 0 && unused != 0x00 {
				return 0, 0, Overflow64Error
			}
			return VarInt64(ret), read, nil
		}
	}
}

func readByte(r io.Reader) (byte, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return 0, err
	}
	return b[0], nil
}
