package shamir

import (
	crand "crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"math/rand/v2"
)

type Shamir struct {
	id     string  // unique identifier to ensure shares were derived from same secret
	field  Gf2m    // field over which to operate
	shares []Share // individual shares to distribute
}

func (shamir Shamir) String() string {
	s := fmt.Sprintf("Secret ID %v\n", shamir.id)
	s += "Shares:\n"
	for n := range shamir.shares {
		s += fmt.Sprintf("  %s\n", shamir.ShareString(n))
	}
	return s
}

func (shamir Shamir) Prefix() string {
	return fmt.Sprintf("shamir-%s-%x", shamir.id, shamir.field.primitivePoly)
}

func (shamir Shamir) ShareLabel(n int) string {
	return fmt.Sprintf("%s-%s", shamir.Prefix(), shamir.shares[n].GetXString())
}

func (shamir Shamir) ShareString(n int) string {
	return fmt.Sprintf("%s-%s", shamir.ShareLabel(n), shamir.shares[n].GetYString())
}

func NewShamirSecret(primitivePoly int, nshares int, threshold int, secret []byte) (*Shamir, error) {

	// input validation
	if threshold > nshares {
		return nil, errors.New("threshold cannot exceed number of shares")
	}
	if (primitivePoly & 0b1) != 1 {
		return nil, errors.New("supplied polynomial cannot be primitive")
	}

	// generate random ID for secret shares
	idbytes := make([]byte, 10)
	if _, err := crand.Read(idbytes); err != nil {
		return nil, errors.New("error reading from random source")
	}

	// initialize the data needed for Shamir's secret sharing scheme
	shamir := &Shamir{
		id:     base32.StdEncoding.EncodeToString(idbytes),
		field:  NewField(primitivePoly),
		shares: make([]Share, nshares),
	}

	// initialize each individual share
	for i := range shamir.shares {
		shamir.shares[i].x = GfElement(i + 1)
		shamir.shares[i].y = make([]GfElement, len(secret))
	}

	// choose new polynomials for each byte in secret
	for i := 0; i < len(secret); i++ {

		// choose random polynomial
		p := make([]GfElement, threshold)
		for i := range p {
			p[i] = GfElement(rand.IntN(shamir.field.GetNelements()))
		}

		// set constant term to be secret
		p[0] = GfElement(secret[i])

		// compute value of polynomial for each of the shares
		for _, share := range shamir.shares {
			share.y[i] = shamir.field.EvaluatePolynomial(p, share.x)
		}
	}

	return shamir, nil
}

func RecoverSecret(primitivePoly int, shares []Share) ([]byte, error) {

	// TODO: check that shares all have same id
	// TODO: check that shares all have different x

	field := NewField(primitivePoly)
	len_secret := len(shares[0].y)
	n_shares := len(shares)

	secret := make([]byte, len_secret)

	x := make([]GfElement, n_shares)
	for s, share := range shares {
		x[s] = share.x
	}

	for i := range len_secret {
		y := make([]GfElement, n_shares)
		for s, share := range shares {
			y[s] = share.y[i]
		}

		// compute L(0) by summing terms l_j(0)
		L := GfElement(0)

		for j := range n_shares {
			ell := GfElement(1)
			for k := range n_shares {
				if k == j {
					continue
				}
				x, err := field.Divide(field.Subtract(GfElement(0), x[k]), field.Subtract(x[j], x[k]))
				if err != nil {
					return nil, err
				}
				ell = field.Multiply(ell, x)
			}
			L = field.Add(L, field.Multiply(y[j], ell))
		}

		secret[i] = byte(L)
	}

	return secret, nil
}