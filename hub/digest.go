package hub

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrNotValidDigest = errors.New("not a valid digest")
)

type Digest struct {
	alg   string
	value string
}

func (d *Digest) String() string {
	return fmt.Sprintf("%s:%s", d.alg, d.value)
}

func (d *Digest) Prefix() string {
	return d.value[0:2]
}

func (d *Digest) Value() string {
	return d.value
}

func ParseDigest(digest string) (*Digest, error) {
	fields := strings.Split(digest, ":")
	if len(fields) != 2 {
		return nil, ErrNotValidDigest
	}
	return &Digest{
		alg:   fields[0],
		value: fields[1],
	}, nil
}

func GetDigest(filePath string) (*Digest, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := sha256.New()

	if _, err = io.Copy(hash, file); err != nil {
		return nil, err
	}

	hashBytes := hash.Sum(nil)

	return &Digest{alg: "sha256", value: fmt.Sprintf("%x", hashBytes)}, nil
}
