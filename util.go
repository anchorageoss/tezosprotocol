package tezosprotocol

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

func serializeBoolean(b bool) byte {
	if b {
		return byte(255)
	}
	return byte(0)
}

func deserializeBoolean(b byte) (bool, error) {
	switch b {
	case 0:
		return false, nil
	case 255:
		return true, nil
	default:
		return false, xerrors.Errorf("byte value %d not a valid boolean encoding", b)
	}
}

func catchOutOfRangeExceptions(r interface{}) error {
	if strings.Contains(fmt.Sprintf("%s", r), "out of range") {
		return xerrors.New("out of bounds exception while parsing operation")
	}
	panic(r)
}
