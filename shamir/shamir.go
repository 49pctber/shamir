package shamir

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	gf "galois_field"
)

var logTable []gf.GfPower
var antilogTable []gf.GfElement

func init() {
	m := 8
	primitivePoly := 0b100011101
	logTable, antilogTable = gf.GenerateGF2mTables(m, primitivePoly)
}

func createShares(secret []byte, nshares, threshold int) error {
	if threshold > nshares {
		return errors.New("threshold cannot exceed number of shares")
	}

	shares := make([][]byte, nshares)
	for i := 0; i < nshares; i++ {
		shares[i] = make([]byte, len(secret))
	}

	for i := 0; i < len(secret); i++ {
		// choose random polynomial
		p, err := chooseRandomPolynomial(secret[i], threshold)
		if err != nil {
			return err
		}

		// compute value of polynomial for each of the shares
		for x := 1; x < nshares+1; x++ {
			y := gf.GfElement(0)
			for d := len(p) - 1; d > -1; d-- {
				// y = y * x + p[d] // what needs to happen, but in GF(256)
				if y == 0 {
					y = gf.GfElement(p[d])
				} else {
					y = antilogTable[(logTable[y]+logTable[x]+255)%255] ^ gf.GfElement(p[d])
				}

			}
			shares[x-1][i] = byte(y)
		}
	}

	// report shares in hex
	for i, share := range shares {
		fmt.Printf("%d: %s\n", i+1, hex.EncodeToString(share))
	}

	return nil

}

func chooseRandomPolynomial(secret byte, threshold int) ([]byte, error) {
	p := make([]byte, threshold)
	_, err := rand.Read(p[1:])
	if err != nil {
		return nil, nil
	}

	p[0] = secret

	return p, nil
}
