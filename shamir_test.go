package shamir

import (
	"bytes"
	"strings"
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

	// reconstruct secret from all available shares
	_, err = RecoverSecret(append(shamir.shares[0:2], shamir.shares[0:2]...))
	if err == nil {
		t.Fatal("should have thrown an error with same shares multiple times")
	}
}

func TestShamirErrors(t *testing.T) {
	share3 := "shamir-7SPFLJYT-11d-3-xYSJU5oTyQcNZHs9SvY"
	share4 := "shamir-7SPFLJYT-11d-4-fu7/+G46PVTx0GBOL5E"
	share4_2 := "shamir-7SPFLJYT-11d-4-fu7/+G46PVTx0GBOL5Efu7/+G46PVTx0GBOL5E"

	want := "This is a test"

	input := strings.Join([]string{share3, share4}, "\n")
	shares, err := NewSharesFromString(input)
	if err != nil {
		t.Errorf("error parsing shares: %v\n", err)
	}

	have, err := RecoverSecret(shares)
	if err != nil {
		t.Errorf("should have reconstructed secret properly: %v\n", err)
	}
	if string(have) != want {
		t.Errorf("error reconstructing secret. Have %s, want %s.", have, want)
	}

	input = strings.Join([]string{share3, share4_2}, "\n")
	shares, err = NewSharesFromString(input)
	if err != nil {
		t.Errorf("error parsing shares: %v\n", err)
	}

	_, err = RecoverSecret(shares)
	if err == nil {
		t.Errorf("should have thrown error\n")
	}
}
