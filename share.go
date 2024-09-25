package shamir

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
)

const SharePrefix string = "shamir"

type Share struct {
	secret_id     string
	primitivePoly int64
	x             GfElement   // x coordinate
	y             []GfElement // y coordinates
}

func NewShare(secret_id string, primitivePoly int64, x GfElement, y []GfElement) Share {
	return Share{secret_id: secret_id, primitivePoly: primitivePoly, x: x, y: y}
}

func NewSharesFromString(input string) ([]Share, error) {
	r := regexp.MustCompile(`shamir-(\w+)-(\w+)-(\w+)-(.+)`)

	shares := make([]Share, 0)
	matches := r.FindAllStringSubmatch(input, -1)
	for _, match := range matches {

		secret_id := string(match[1])

		primitivePoly, err := strconv.ParseInt(string(match[2]), 16, 64)
		if err != nil {
			return nil, err
		}

		xdata, err := strconv.ParseInt(string(match[3]), 10, 64)
		if err != nil {
			return nil, err
		}

		ydata, err := base64.RawStdEncoding.DecodeString(string(match[4]))
		if err != nil {
			return nil, err
		}

		x := GfElement(xdata)
		y := make([]GfElement, len(ydata))
		for i := range ydata {
			y[i] = GfElement(ydata[i])
		}

		shares = append(shares, NewShare(secret_id, primitivePoly, x, y))
	}

	return shares, nil

}

func (share Share) ShareLabel() string {
	return fmt.Sprintf("%s-%s-%x-%s", SharePrefix, share.secret_id, share.primitivePoly, share.GetXString())
}

func (share Share) String() string {
	return fmt.Sprintf("%s-%s", share.ShareLabel(), share.GetYString())
}

func (share Share) GetSecretId() string {
	return share.secret_id
}

func (share Share) GetPrimitivePoly() int64 {
	return share.primitivePoly
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
