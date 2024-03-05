package repository

import (
	"io"
	"mime/multipart"
	"os"
	"path"
)

type Store struct{}

func NewStoreRepository() *Store {
	return &Store{}
}

func (s *Store) FolderCreate(name string) error {
	newPath := path.Join("uploads", name)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SaveFile(file *multipart.File, fileName string) (string, error) {
	err := s.FolderCreate(fileName)
	if err != nil {
		return "", err
	}
	filePath := path.Join("uploads", fileName, "original")
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, *file)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
