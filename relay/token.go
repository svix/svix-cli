package relay

import "github.com/segmentio/ksuid"

func GenerateToken() string {
	return ksuid.New().String()
}
