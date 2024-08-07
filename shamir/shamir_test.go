package shamir

import "testing"

func TestShamir(t *testing.T) {
	createShares([]byte("bryan"), 2, 2)
}
