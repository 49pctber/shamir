package galoisfield

import (
	"math"
)

type GfElement int
type GfPower int

func GenerateGF2mTables(m int, primitivePoly int) ([]GfPower, []GfElement) {

	const q = 2 // will only produce GF(2^m)

	n_elements := 1
	for i := 0; i < m; i++ {
		n_elements *= q
	}

	log := make([]GfPower, n_elements)       // polynomial to power representation
	antilog := make([]GfElement, n_elements) // power to polynomial representation

	var poly GfElement = 1

	for power := GfPower(0); int(power) < n_elements-1; power++ {
		antilog[power] = poly
		log[poly] = power

		poly <<= 1
		if poly&(0b1<<m) != 0 {
			poly ^= GfElement(primitivePoly)
		}
	}

	log[1] = 0
	log[0] = math.MinInt

	return log, antilog
}
