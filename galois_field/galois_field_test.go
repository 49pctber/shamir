package galoisfield

import (
	"testing"
)

func TestTables2(t *testing.T) {
	m := 2
	primitivePoly := 0b111
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d", element, polyBack)
		}
	}
}

func TestTables3(t *testing.T) {
	m := 3
	primitivePoly := 0b1011
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d", element, polyBack)
		}
	}
}

func TestTables4(t *testing.T) {
	m := 4
	primitivePoly := 0b10011
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d", element, polyBack)
		}
	}

	if have, want := field.Add(3, 6), GfElement(5); have != want {
		t.Errorf("3+6=5, not %d", have)
	}

	if have, want := field.Add(15, 15), GfElement(0); have != want {
		t.Errorf("15+15=0, not %d", have)
	}

	if have, want := field.Multiply(3, 6), GfElement(10); have != want {
		t.Errorf("3*6=10, not %d", have)
	}

	if have, want := field.Multiply(2, 9), GfElement(1); have != want {
		t.Errorf("2*9=1, not %d", have)
	}

	if have, want := field.Multiply(0, 10), GfElement(0); have != want {
		t.Errorf("0*10=0, not %d", have)
	}

	for i := 1; i < 16; i++ {
		want := GfElement(1)
		have, err := field.Divide(GfElement(i), GfElement(i))
		if err != nil {
			t.Error(err)
		}
		if have != want {
			t.Errorf("should be identity, not %d", have)
		}
	}

	for i := 0; i < 16; i++ {
		want := GfElement(i)
		have, err := field.Divide(GfElement(i), GfElement(1))
		if err != nil {
			t.Error(err)
		}
		if have != want {
			t.Errorf("should be original number %d, not %d", want, have)
		}
	}

	_, err := field.Divide(GfElement(1), GfElement(0))
	if err == nil {
		t.Error("should have thrown a division by zero error")
	}

	for j := 0; j < 16; j++ {
		for i := 1; i < 16; i++ {
			want := GfElement(j)
			have, err := field.Divide(field.Multiply(want, GfElement(i)), GfElement(i))
			if err != nil {
				t.Error(err)
			}
			if have != want {
				t.Errorf("%d*%d/%d=%d, not %d", want, i, i, want, have)
			}
		}
	}

}

func TestTables5_1(t *testing.T) {
	m := 5
	primitivePoly := 0b100101
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d", element, polyBack)
		}
	}
}

func TestTables5_2(t *testing.T) {
	m := 5
	primitivePoly := 0b110111
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d", element, polyBack)
		}
	}
}

func TestTables8_1(t *testing.T) {
	m := 8
	primitivePoly := 0b100011101
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", element, polyBack, element, polyBack)
		}
	}
}

func TestTables11_1(t *testing.T) {
	m := 11
	primitivePoly := 0b100000000101
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", element, polyBack, element, polyBack)
		}
	}
}

func TestTables8_2(t *testing.T) {
	m := 8
	primitivePoly := 0b101011111 // different primitive polynomial
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", element, polyBack, element, polyBack)
		}
	}
}

func TestTables21(t *testing.T) {
	m := 21
	primitivePoly := 0b1000000000000000000101
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", element, polyBack, element, polyBack)
		}
	}
}

func TestTables23_1(t *testing.T) {
	m := 23
	primitivePoly := 0b100000000000000000100001
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", element, polyBack, element, polyBack)
		}
	}
}

func TestTables23_2(t *testing.T) {
	m := 23
	primitivePoly := 0b101000000000000010100001
	field := NewField(m, primitivePoly)

	for element := GfElement(1); element < 1<<m; element++ {
		power := field.logTable[element]
		polyBack := field.antilogTable[power]
		if element != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", element, polyBack, element, polyBack)
		}
	}
}

func TestEvaluatePolynomial(t *testing.T) {

	m := 8
	primitivePoly := 0b100011101
	field := NewField(m, primitivePoly)

	if have, want := field.EvaluatePolynomial([]GfElement{32}, 0), GfElement(32); have != want {
		t.Errorf("polynomial is constant 32")
	}

	if have, want := field.EvaluatePolynomial([]GfElement{32}, 17), GfElement(32); have != want {
		t.Errorf("polynomial is constant and should not vary with x")
	}

	if have, want := field.EvaluatePolynomial([]GfElement{32}, 255), GfElement(32); have != want {
		t.Errorf("polynomial is constant and should not vary with x")
	}

	if have, want := field.EvaluatePolynomial([]GfElement{3, 5}, 0), GfElement(3); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{3, 5}, 1), GfElement(6); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{3, 5}, 2), GfElement(9); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{3, 5}, 3), GfElement(12); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{3, 5}, 4), GfElement(23); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{4, 2, 1}, 5), GfElement(31); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{241, 170}, 27), GfElement(89); have != want {
		t.Errorf("have %d, want %d", have, want)
	}

	if have, want := field.EvaluatePolynomial([]GfElement{241, 170, 180}, 27), GfElement(25); have != want {
		t.Errorf("have %d, want %d", have, want)
	}
}
