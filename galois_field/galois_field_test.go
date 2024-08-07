package galoisfield

import (
	"fmt"
	"testing"
)

func TestTables2(t *testing.T) {
	m := 2
	primitivePoly := 0b111
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d", poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables3(t *testing.T) {
	m := 3
	primitivePoly := 0b1011
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d", poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables4(t *testing.T) {
	m := 4
	primitivePoly := 0b10011
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d", poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables5_1(t *testing.T) {
	m := 5
	primitivePoly := 0b100101
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d", poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables5_2(t *testing.T) {
	m := 5
	primitivePoly := 0b110111
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d", poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables8_1(t *testing.T) {
	m := 8
	primitivePoly := 0b100011101
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", poly, polyBack, poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables11_1(t *testing.T) {
	m := 11
	primitivePoly := 0b100000000101
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", poly, polyBack, poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables8_2(t *testing.T) {
	m := 8
	primitivePoly := 0b101011111 // different primitive polynomial
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", poly, polyBack, poly, polyBack)
		}
	}

	if t.Failed() {
		fmt.Println(logTable)
	}
}

func TestTables21(t *testing.T) {
	m := 21
	primitivePoly := 0b1000000000000000000101
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", poly, polyBack, poly, polyBack)
		}
	}
}

func TestTables23_1(t *testing.T) {
	m := 23
	primitivePoly := 0b100000000000000000100001
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", poly, polyBack, poly, polyBack)
		}
	}
}

func TestTables23_2(t *testing.T) {
	m := 23
	primitivePoly := 0b101000000000000010100001
	logTable, antilogTable := GenerateGF2mTables(m, primitivePoly)

	for poly := GfElement(1); poly < 1<<m; poly++ {
		power := logTable[poly]
		polyBack := antilogTable[power]
		if poly != polyBack {
			t.Errorf("exp^log(%d) != %d (expected %d but got %d)", poly, polyBack, poly, polyBack)
		}
	}
}
