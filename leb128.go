package leb128

import (
	"errors"
)

var sevenbits = [...]byte{
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
	0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
	0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
	0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
}

type bs []byte

// SLeb128 : There are 2 versions of LEB128: unsigned LEB128 and signed128.
// The decoder must know whether the encoded value is unsigned or signed.
type SLeb128 struct {
	bs
}

// ULeb128 : Unsigned LEB128
type ULeb128 struct {
	bs
}

// AppendSLeb128 : Append two `Sleb128` into one `SLeb128`
func AppendSLeb128(x, y SLeb128) SLeb128 {
	return SLeb128{append(x.bs, y.bs...)}
}

// AppendULeb128 : Append two `Uleb128` into one `ULeb128`
func AppendULeb128(x, y ULeb128) ULeb128 {
	return ULeb128{append(x.bs, y.bs...)}
}

// EncodeFromUint64 : Encode `uint64` typed integer into `ULeb128`
func EncodeFromUint64(input uint64) ULeb128 {
	// https://golang.org/src/cmd/internal/dwarf/dwarf.go
	if input < 0x80 {
		return ULeb128{[]byte{sevenbits[input]}}
	}

	// https://en.wikipedia.org/wiki/LEB128
	var bytes []byte
	for {
		b := uint8(input & 0x7f)
		input = input >> 7

		if input != 0 {
			b |= 0x80
		}

		bytes = append(bytes, b)
		if b&0x80 == 0 {
			break
		}
	}
	return ULeb128{bytes}
}

// EncodeFromInt64 : Encode `int64` typed integer into `SLeb128`
func EncodeFromInt64(input int64) SLeb128 {
	if input >= 0 && input <= 0x3f {
		return SLeb128{[]byte{sevenbits[input]}}
	} else if input < 0 && input >= ^0x3f {
		return SLeb128{[]byte{sevenbits[0x80+input]}}
	}

	var bytes []byte
	for {
		b := uint8(input & 0x7f)
		signBit := uint8(input & 0x40) // sign bit is second high order bit
		input = input >> 7

		if (input != 0 || signBit != 0) && (input != -1 || signBit == 0) {
			b |= 0x80 // set high order bit of byte
		}

		bytes = append(bytes, b)
		if b&0x80 == 0 {
			break
		}
	}
	return SLeb128{bytes}
}

// DecodeULeb128 :
func DecodeULeb128(src ULeb128) (result uint64, err error) {
	length := uint8(len(src.bs) & 0xff)
	if length > 10 {
		length = 10
	}

	var i uint8
	for i = 0; i < length; i++ {
		result |= uint64(src.bs[i]&0x7f) << (7 * i) // shift
		if src.bs[i]&0x80 == 0 {
			if uint8(i+1) == 0 {
				result = 0
				err = errors.New("bad encoding")
				return
			}
		}

	}
	return
}

// DecodeSLeb128 :
func DecodeSLeb128(src SLeb128) (result int64, err error) {
	length := uint8(len(src.bs) & 0xff)
	if length > 10 {
		length = 10
	}

	var i uint8
	for i = 0; i < length; i++ {
		result |= int64(src.bs[i]&0x7f) << (7 * i)
		if src.bs[i]&0x80 == 0 {
			if src.bs[i]&0x40 != 0 { // bitflag
				result |= ^0 << (7 * (i + 1))
			}
			if uint8(i+1) == 0 {
				result = 0
				err = errors.New("bad encoding")
				return
			}
		}
	}
	return
}
