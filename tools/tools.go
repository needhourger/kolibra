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

type TxtReader struct {
	F      *os.File
	Reader *bufio.Reader
}

func OpenTxtByEncode(path string) (*TxtReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(f)
	dumpedBytes, err := reader.Peek(1024)
	f.Seek(0, io.SeekStart)
	reader.Reset(f)
	if err != nil {
		return nil, err
	}
	_, name, _ := charset.DetermineEncoding(dumpedBytes, "text/plain")
	log.Printf("Detected encoding: %s", name)
	if name == "utf-8" {
		return &TxtReader{f, reader}, nil
	}
	newDecodedReader := transform.NewReader(f, simplifiedchinese.GBK.NewDecoder().Transformer)
	return &TxtReader{f, bufio.NewReader(newDecodedReader)}, nil
}
