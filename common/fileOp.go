package common

import (
	"errors"
	"os"
	"path/filepath"
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

func (s FileOp) DirCreate(path, name string) error {
	return os.MkdirAll(filepath.Join(s.BasePath, path, name), 640)
}

func (s FileOp) FileCreate(name string, b []byte) (int, error) {
	file, err := os.OpenFile(filepath.Join(s.BasePath, name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	return file.Write(b)
}

func (s FileOp) FileRemoveAll(path string) error {
	return os.RemoveAll(filepath.Join(s.BasePath, path))
}

func (s FileOp) FileRename(oldpath, oldid, newpath, newid string) error {
	return os.Rename(filepath.Join(s.BasePath, oldpath, oldid), filepath.Join(s.BasePath, newpath, newid))
}
