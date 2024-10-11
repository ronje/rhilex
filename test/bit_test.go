package test

import (
	"testing"

	"github.com/hootrhino/rhilex/rhilexlib"
)

func TestReverseBitOrder(t *testing.T) {
	t.Log(t, byte(0b1011_1111), rhilexlib.ReverseBits(byte(0b1111_1101)))
	t.Log(t, byte(0b1100_0000), rhilexlib.ReverseBits(byte(0b0000_0011)))
	t.Log(t, byte(0b0000_0101), rhilexlib.ReverseBits(byte(0b1010_0000)))
	t.Log(t, byte(0b1010_1010), rhilexlib.ReverseBits(byte(0b0101_0101)))
	t.Log(t, []byte{3, 2, 1}, rhilexlib.ReverseByteOrder([]byte{1, 2, 3}))
}
