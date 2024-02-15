package local

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
)

type LocalStore struct{}

func NewLocalStore() *LocalStore {
	return &LocalStore{}
}

func (ls *LocalStore) Ls(filepath string) (map[string]string, map[string]string, error) {
	folders := make(map[string]string, 0)
	files := make(map[string]string, 0)

	entries, err := os.ReadDir(path.Join("./uploads", filepath))
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	for _, e := range entries {
		if e.IsDir() {
			folders[e.Name()] = path.Join(filepath, e.Name())
		} else {
			files[e.Name()] = path.Join(filepath, e.Name())
		}
	}

	return folders, files, nil
}

func (ls *LocalStore) FolderCreate(name string) error {
	newPath := path.Join("uploads", name)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (ls *LocalStore) Create(file *multipart.File, filepath string, filename string) error {
	newPath := path.Join("uploads", filepath)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return err
	}
	dst, err := os.Create(path.Join("uploads", filepath, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, *file); err != nil {
		return err
	}
	return nil
}
