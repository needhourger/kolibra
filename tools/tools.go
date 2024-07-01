package tools

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"regexp"

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
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

func GetFileEncoded(path string) (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	rawBytes, err := reader.Peek(4096)
	if err != nil {
		return "", err
	}
	_, name, _ := charset.DetermineEncoding(rawBytes, "")
	return name, nil
}

func Gbk2utf8String(rawBytes []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(rawBytes), simplifiedchinese.GBK.NewDecoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(d), nil
}

func IsMatchString(str string, pattern string) (bool, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	match := r.MatchString(str)
	return match, nil
}
