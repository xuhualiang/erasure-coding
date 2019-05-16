package main

import (
	"fmt"
	. "org/coding/ec"
)

type MultTable [256][256]uint8
type InvTable  [256]uint8

func (mt *MultTable) Print() {
	fmt.Printf("var MT = [][]uint8{")
	for i := 0; i < len(mt); i += 1{
		if i != 0 {
			fmt.Printf(",")
		}
		fmt.Print("\n\t{")

		for j := 0; j < len(mt[i]); j += 1{
			if j != 0 {
				fmt.Printf(", ")
			}
			if j != 0 && j % 16 == 0 {
				fmt.Printf("\n\t ")
			}
			fmt.Printf("0x%02x", mt[i][j])
		}
		fmt.Print("}")
	}
	fmt.Printf("}\n\n");
}

func (mt *MultTable) Verify()  {
	assert(mt[0][0] == GF_ZERO)

	for a := 1; a < 256; a += 1 {
		// a * 0 = 0
		assert(mt[a][GF_ZERO] == GF_ZERO)
		// 0 * a = 0
		assert(mt[GF_ZERO][a] == GF_ZERO)
		// a * 1 = a
		assert(mt[a][GF_ONE] == uint8(a))
		// 1 * a = a
		assert(mt[GF_ONE][a] == uint8(a))

		for b := 0; b < 256; b += 1 {
			// a * b = b * a
			assert(mt[a][b] == mt[b][a])

			for c := 0; c < 256; c += 1 {
				// (a * b) * c == a * (b * c)
				assert(mt[mt[a][b]][c] == mt[a][mt[b][c]])

				// (a + b) * c == a * c + b * c
				assert(mt[a ^ b][c] == mt[a][c] ^ mt[b][c])
			}
		}
	}
}

func (invt *InvTable) Verify(mt *MultTable)  {
	// 1/0 invalid
	assert(invt[0] == 0)

	for a := 1; a < 256; a += 1 {
		// a * 1/a == 1
		assert(mt[a][invt[a]] == GF_ONE)
	}
}

func (invt *InvTable) Print() {
	fmt.Printf("var INVT = []uint8{")
	for i := 0; i < len(invt); i += 1 {
		if i != 0 {
			fmt.Printf(", ")
		}
		if i % 16 == 0 {
			fmt.Printf("\n\t")
		}
		fmt.Printf("0x%02x", invt[i])
	}
	fmt.Printf("}\n\n")
}

func gfMult(a, b uint8) uint8 {
	r := uint8(0)

	for b != 0 {
		if b & 0x01 == 0x01 {
			r ^= a
		}

		if a & 0x80 == 0x80 {
			a = (a << 1) ^ POLY
		} else {
			a = a << 1
		}

		b >>= 1
	}
	return r
}

func genMultTable() *MultTable {
	mt := &MultTable{}

	for a := 0; a < 256; a += 1 {
		for b := 0; b < 256; b += 1 {
			mt[a][b] = gfMult(uint8(a), uint8(b))
		}
	}
	return mt
}

func genInvTable() *InvTable {
	invt := &InvTable{}

	for a := 1; a < 256; a += 1 {
		for b := 1; b < 256; b += 1 {
			if gfMult(uint8(a), uint8(b)) == GF_ONE {
				invt[a] = uint8(b)
				break
			}
		}
	}

	return invt
}

func assert(cond bool)  {
	if !cond {
		panic("assert failure")
	}
}

func main() {
	mult := genMultTable()
	invt := genInvTable()

	mult.Verify()
	invt.Verify(mult)

	fmt.Println("package ec\n\n")

	mult.Print()
	invt.Print()

}
