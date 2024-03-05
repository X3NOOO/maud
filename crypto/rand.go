package crypto

import (
	"crypto/rand"
	"fmt"
)

func RandomString(length uint) (string, error) {
	b := make([]byte, length+2)
	_, err := rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2], err
}
