package shamir

import (
	"fmt"
	"testing"
)

func TestShamir(t *testing.T) {
	shamir, err := NewShamir(0b100011101, 5, 2, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(shamir)
}

func TestComputeDegree(t *testing.T) {
	if have, want := ComputeDegree(0b1), 0; have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := ComputeDegree(0b111), 2; have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := ComputeDegree(0b100011101), 8; have != want {
		t.Errorf("have %d, want %d", have, want)
	}
}
