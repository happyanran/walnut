package common

import (
	"errors"
	"os"
)

type FileOp struct {
	BasePath string
}

func NewFileOp(s ServerConf) (*FileOp, error) {
	if s.Data == "" {
		return nil, errors.New("Data dir should not be null.")
	}

	return &FileOp{BasePath: s.Data}, os.MkdirAll(s.Data, 640)
}

func (s FileOp) DirCreate(path string) error {
	return os.MkdirAll(path, 640)
}
