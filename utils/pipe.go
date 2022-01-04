package utils

import (
	"io"
	"os"
	"strings"
)

func IsStdinReadable() (bool, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}
	return (info.Mode() & os.ModeCharDevice) == 0, nil
}

func ReadStdin() ([]byte, error) {
	isReadable, err := IsStdinReadable()
	if err != nil {
		return nil, err
	}
	if isReadable {
		in, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		return []byte(strings.Trim(string(in), " \r\n")), nil
	}
	return []byte{}, nil
}

func IsTTY(f *os.File) (isTTY bool, hasBytes bool, err error) {
	info, err := f.Stat()
	if err != nil {
		return false, false, err
	}

	isTTY = info.Mode()&os.ModeCharDevice != 0
	hasBytes = info.Size() > 0
	return isTTY, hasBytes, nil
}
