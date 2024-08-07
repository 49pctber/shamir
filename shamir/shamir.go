package shamir

import (
	crand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	gf "galois_field"
	"math/rand/v2"
)

type Share struct {
	id string         // id associated with a given secret
	x  gf.GfElement   // x coordinate
	y  []gf.GfElement // y coordinates
}

func (share Share) String() string {
	return fmt.Sprintf("Share #%d: %v", share.x, share.y)
}

type Shamir struct {
	id     string  // unique identifier to ensure shares were derived from same secret
	field  gf.Gf2m // field over which to operate
	shares []Share // individual shares to distribute
}

func (shamir Shamir) String() string {
	s := fmt.Sprintf("Secret %v\n", shamir.id)
	s += fmt.Sprintf("  Field Parameters: %v\n", shamir.field)
	s += "  Shares:\n"
	for _, share := range shamir.shares {
		s += fmt.Sprintf("    %v\n", share)
	}
	return s
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

func NewShamir(primitivePoly int, nshares int, threshold int, secret []byte) (*Shamir, error) {

	// input validation
	if threshold > nshares {
		return nil, errors.New("threshold cannot exceed number of shares")
	}

	if (primitivePoly & 0b1) != 1 {
		return nil, errors.New("supplied polynomial cannot be primitive")
	}

	// compute degree of primitive polynomial
	m := ComputeDegree(primitivePoly)

	// generate random ID for secret shares
	idbytes := make([]byte, 15)
	if _, err := crand.Read(idbytes); err != nil {
		return nil, errors.New("error reading from random source")
	}

	// initialize the data needed for Shamir's secret sharing scheme
	shamir := &Shamir{
		id:     base64.StdEncoding.EncodeToString(idbytes),
		field:  gf.NewField(m, primitivePoly),
		shares: make([]Share, nshares),
	}

	// initialize each individual share
	for i := range shamir.shares {
		shamir.shares[i].id = shamir.id
		shamir.shares[i].x = gf.GfElement(i + 1)
		shamir.shares[i].y = make([]gf.GfElement, len(secret))
	}

	// choose new polynomials for each byte in secret
	for i := 0; i < len(secret); i++ {

		// choose random polynomial
		p := make([]gf.GfElement, threshold)
		for i := range p {
			p[i] = gf.GfElement(rand.IntN(shamir.field.GetNelements()))
		}

		// set constant term to be secret
		p[0] = gf.GfElement(secret[i])

		// compute value of polynomial for each of the shares
		for _, share := range shamir.shares {
			share.y[i] = shamir.field.EvaluatePolynomial(p, share.x)
		}
	}

	return shamir, nil
}
