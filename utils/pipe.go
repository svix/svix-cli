package utils

import (
	"os"
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
