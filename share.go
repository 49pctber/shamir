package shamir

import (
	"encoding/base64"
	"fmt"
)

type Share struct {
	x GfElement   // x coordinate
	y []GfElement // y coordinates
}

func NewShare(x GfElement, y []GfElement) Share {
	return Share{x: x, y: y}
}

func (share Share) String() string {
	return fmt.Sprintf("%s-%s", share.GetXString(), share.GetYString())
}

func (share Share) GetXString() string {
	return fmt.Sprintf("%d", share.x)
}

func (share Share) GetYString() string {
	b := make([]byte, len(share.y))
	for i := range b {
		b[i] = byte(share.y[i])
	}
	return base64.RawStdEncoding.EncodeToString(b)
}
