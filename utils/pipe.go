package utils

import (
	"io"
	"os"
	"strings"
)

func IsPipeWithData(f *os.File) (bool, error) {
	info, err := f.Stat()
	if err != nil {
		return false, err
	}

	if info.Mode()&os.ModeCharDevice != 0 || // stdin is not a Unix character device
		info.Size() <= 0 { // stdin has no bytes
		return false, nil
	}
	return true, nil
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
