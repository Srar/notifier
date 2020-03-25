package network

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"strings"
)

func DecodeString(p []byte, len int, desc *string) []byte {
	s := string(p[0:len])
	s = strings.Trim(s, "\000")
	*desc = strings.TrimSpace(s)
	return p[len:]
}

func DecodeBytes(p []byte, len int, desc []byte) []byte {
	desc = p[0:len]
	return p[len:]
}

func DecodeHex(p []byte, len int, desc *string) []byte {
	s := fmt.Sprintf("%x", p[0:len])
	*desc = s
	return p[len:]
}

func EncodeBool(p []byte, c bool) []byte {
	bit := uint8(0)
	if c {
		bit = 1
	}
	return Encode8u(p, bit)
}

/* encode 8 bits unsigned int */
func Encode8u(p []byte, c byte) []byte {
	p[0] = c
	return p[1:]
}

func DecodeBool(p []byte, c *bool) []byte {
	var a uint8
	var b = false
	afterDecodePayload := Decode8u(p, &a)
	if a > 0 {
		b = true
	}
	c = &b
	return afterDecodePayload
}

/* decode 8 bits unsigned int */
func Decode8u(p []byte, c *byte) []byte {
	*c = p[0]
	return p[1:]
}

/* encode 16 bits unsigned int (lsb) */
func Encode16u(p []byte, w uint16) []byte {
	binary.LittleEndian.PutUint16(p, w)
	return p[2:]
}

/* decode 16 bits unsigned int (lsb) */
func Decode16u(p []byte, w *uint16) []byte {
	*w = binary.LittleEndian.Uint16(p)
	return p[2:]
}

/* encode 32 bits unsigned int (lsb) */
func Encode32u(p []byte, l uint32) []byte {
	binary.LittleEndian.PutUint32(p, l)
	return p[4:]
}

/* decode 32 bits unsigned int (lsb) */
func Decode32u(p []byte, l *uint32) []byte {
	*l = binary.LittleEndian.Uint32(p)
	return p[4:]
}

/* encode 64 bits unsigned int (lsb) */
func Encode64u(p []byte, l uint64) []byte {
	binary.LittleEndian.PutUint64(p, l)
	return p[4:]
}

/* decode 64 bits unsigned int (lsb) */
func Decode64u(p []byte, l *uint64) []byte {
	*l = binary.LittleEndian.Uint64(p)
	return p[8:]
}

func EncodeFloat64(p []byte, l float64) []byte {
	binary.LittleEndian.PutUint64(p, math.Float64bits(l))
	return p[8:]
}

func DecodeFloat64(p []byte, l *float64) []byte {
	*l = math.Float64frombits(binary.LittleEndian.Uint64(p))
	return p[8:]
}

func Min(a, b uint32) uint32 {
	if a <= b {
		return a
	}
	return b
}

func MinUint16(a, b uint16) uint16 {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b uint32) uint32 {
	if a >= b {
		return a
	}
	return b
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func Bound(lower, middle, upper uint32) uint32 {
	return Min(Max(lower, middle), upper)
}

func Diff(later, earlier uint32) int32 {
	return (int32)(later - earlier)
}
