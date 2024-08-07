package galoisfield

import (
	"errors"
	"fmt"
	"math"
)

type GfElement int
type GfPower int

type Gf2m struct {
	m             int
	n_elements    int
	primitivePoly int
	logTable      []GfPower
	antilogTable  []GfElement
}

func (field Gf2m) String() string {
	return fmt.Sprintf("GF(2^%d) using primitive polynomial 0x%x", field.m, field.primitivePoly)
}

func (field Gf2m) GetNelements() int {
	return field.n_elements
}

// computes the degree of a given polynomial
func ComputeDegree(poly int) int {
	m := 0
	for poly > 1 {
		poly >>= 1
		m += 1
	}
	return m
}

func NewField(primitivePoly int) Gf2m {

	const q = 2 // will only produce GF(2^m)
	m := ComputeDegree(primitivePoly)

	n_elements := 1
	for i := 0; i < m; i++ {
		n_elements *= q
	}

	lut := Gf2m{
		m:             m,
		n_elements:    n_elements,
		primitivePoly: primitivePoly,
		logTable:      make([]GfPower, n_elements),
		antilogTable:  make([]GfElement, n_elements),
	}

	var poly GfElement = 1

	for power := GfPower(0); int(power) < n_elements-1; power++ {
		lut.antilogTable[power] = poly
		lut.logTable[poly] = power

		poly <<= 1
		if poly&(0b1<<m) != 0 {
			poly ^= GfElement(primitivePoly)
		}
	}

	lut.logTable[0] = math.MinInt
	lut.logTable[1] = 0
	lut.antilogTable[len(lut.antilogTable)-1] = lut.antilogTable[0]

	return lut
}

// add two elements in the field
func (field Gf2m) Add(a, b GfElement) GfElement {
	return a ^ b
}

// subtract two elements in the field
// addition and subtraction are the same in GF(2^m)
func (field Gf2m) Subtract(a, b GfElement) GfElement {
	return a ^ b
}

// multiply two elements in the field
func (field Gf2m) Multiply(a, b GfElement) GfElement {

	if a == 0 || b == 0 {
		return GfElement(0)
	}

	loga := field.logTable[a]
	logb := field.logTable[b]
	logc := (loga + logb) % GfPower(field.n_elements-1)
	return field.antilogTable[logc]
}

// divide a by b
func (field Gf2m) Divide(a, b GfElement) (GfElement, error) {

	if b == 0 {
		return GfElement(0), errors.New("division by zero")
	}

	if a == 0 {
		return GfElement(0), nil
	}

	loga := field.logTable[a]
	logb := field.logTable[b]
	logc := (loga + GfPower(field.n_elements-1) - logb) % GfPower(field.n_elements-1)
	return field.antilogTable[logc], nil
}

// evaluate a polynomial over a field
func (field Gf2m) EvaluatePolynomial(p []GfElement, x GfElement) (y GfElement) {
	y = 0
	for d := len(p) - 1; d > -1; d-- {
		y = field.Multiply(y, GfElement(x))
		y = field.Add(y, GfElement(p[d]))
	}
	return y
}
