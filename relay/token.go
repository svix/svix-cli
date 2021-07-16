package relay

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateToken() (string, error) {
	return gonanoid.New(27)
}
