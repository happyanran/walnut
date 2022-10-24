package common

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

type FileOp struct {
	BasePath string
}

var fileOpLog *logrus.Logger

func NewFileOp(s ServerConf, log *logrus.Logger) (*FileOp, error) {
	fileOpLog = log

	if s.Data == "" {
		return nil, errors.New("Data dir should not be null.")
	}

	return &FileOp{BasePath: s.Data}, os.MkdirAll(s.Data, 640)
}

func (s FileOp) DirCreate(path string) error {
	return os.MkdirAll(filepath.Join(s.BasePath, path), 640)
}

func (s FileOp) FileUpload(path, name string, f func(string) error) error {
	if err := s.DirCreate(path); err != nil {
		return err
	}

	return f(filepath.Join(s.BasePath, path, name))
}

func (s FileOp) ImgResize(path, srcName, destName string, width, height int) {
	src, err := imaging.Open(filepath.Join(s.BasePath, path, srcName))
	if err != nil {
		fileOpLog.Error(err)
	}

	des := imaging.Resize(src, width, height, imaging.Lanczos)

	err = imaging.Save(des, filepath.Join(s.BasePath, path, destName))
	if err != nil {
		fileOpLog.Error(err)
	}
}

func (s FileOp) FileCreate(name string, b []byte) (int, error) {
	file, err := os.OpenFile(filepath.Join(s.BasePath, name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return 0, err
	}

	defer file.Close()

	return file.Write(b)
}
