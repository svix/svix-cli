package version

import (
	"fmt"
)

var Version = "source"

func String() string {
	return fmt.Sprintf("svix version: %s\n", Version)
}
