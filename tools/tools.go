package tools

import (
	"bufio"
	"crypto/sha256"
	"io"
	"log"
	"os"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

func OpenFile(path string) (*os.File, *bufio.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	reader := bufio.NewReader(f)
	dumpedBytes, err := reader.Peek(1024)
	if err != nil {
		return nil, nil, err
	}
	_, name, _ := charset.DetermineEncoding(dumpedBytes, "text/plain")
	log.Printf("Detected encoding: %s", name)
	if name == "utf-8" {
		return f, reader, nil
	}
	newDecodedReader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder().Transformer)
	return f, bufio.NewReader(newDecodedReader), nil
}
