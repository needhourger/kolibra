package tools

import (
	"crypto/sha256"
	"io"
	"os"
)

func CalculateFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	ret := hash.Sum(nil)
	return string(ret), nil
}
