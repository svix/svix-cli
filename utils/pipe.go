package utils

import (
	"encoding/json"
	"io"
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

func TryMarshallPipe(v interface{}) error {
	isPipe, err := IsPipeWithData(os.Stdin)
	if err != nil {
		return err
	}
	if isPipe {
		in, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(in, v); err != nil {
			return err
		}
	}
	return nil
}
