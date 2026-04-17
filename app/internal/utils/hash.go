package utils

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func Hash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func CompareHash(s string, hash string) int {
	return strings.Compare(Hash(s), hash)
}
