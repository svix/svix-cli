package utils

import (
	"io"
	"os"
	"strings"
)

func IsPipeWithData(f *os.File) (bool, error) {
	isTTY, hasData, err := IsTTY(f)
	if err != nil {
		return false, err
	}
	return !isTTY && hasData, nil
}

func ReadPipe() ([]byte, error) {
	isPipe, err := IsPipeWithData(os.Stdin)
	if err != nil {
		return nil, err
	}
	if isPipe {
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
