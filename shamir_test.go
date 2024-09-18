package shamir

import (
	"bytes"
	"testing"
)

func TestShamir(t *testing.T) {
	// parameters
	secret := []byte("This is a secret 🤫")
	primitivePoly := 0x11d
	nshares := 6
	threshold := 2

	// compute shares
	shamir, err := NewShamirSecret(primitivePoly, nshares, threshold, secret)
	if err != nil {
		t.Fatal(err)
	}

	// reconstruct secret from shares
	recovered_secret, err := RecoverSecret(shamir.shares[0:2])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if !bytes.Equal(secret, recovered_secret) {
		t.Fatalf("have %v, want %v", recovered_secret, secret)
	}

	// reconstruct secret from different shares
	recovered_secret, err = RecoverSecret(shamir.shares[2:4])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if !bytes.Equal(secret, recovered_secret) {
		t.Fatalf("have %v, want %v", recovered_secret, secret)
	}

	// reconstruct secret from different shares
	recovered_secret, err = RecoverSecret(shamir.shares[4:6])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if !bytes.Equal(secret, recovered_secret) {
		t.Fatalf("have %v, want %v", recovered_secret, secret)
	}
}

func TestShamir_2(t *testing.T) {
	// parameters
	secret := []byte("You just lost the game.")
	primitivePoly := 0x11d
	nshares := 6
	threshold := 4

	// compute shares
	shamir, err := NewShamirSecret(primitivePoly, nshares, threshold, secret)
	if err != nil {
		t.Fatal(err)
	}

	// should not be able to reconstruct secret from 2 shares
	recovered_secret, err := RecoverSecret(shamir.shares[0:2])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if bytes.Equal(secret, recovered_secret) {
		t.Fatal("you cheated somehow...you shouldn't be able to reconstruct the secret")
	}

	// should not be able to reconstruct secret from 2 shares
	recovered_secret, err = RecoverSecret(shamir.shares[0:3])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if bytes.Equal(secret, recovered_secret) {
		t.Fatal("you cheated somehow...you shouldn't be able to reconstruct the secret")
	}

	// reconstruct secret from minimum number of shares
	recovered_secret, err = RecoverSecret(shamir.shares[0:4])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if !bytes.Equal(secret, recovered_secret) {
		t.Fatalf("have %v, want %v", recovered_secret, secret)
	}

	// reconstruct secret from all available shares
	recovered_secret, err = RecoverSecret(shamir.shares[0:6])
	if err != nil {
		t.Fatal(err)
	}

	// check that everything went well
	if !bytes.Equal(secret, recovered_secret) {
		t.Fatalf("have %v, want %v", recovered_secret, secret)
	}
}
