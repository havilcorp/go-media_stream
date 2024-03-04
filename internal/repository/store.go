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

// func (s *Store) Ls(filepath string) (map[string]string, map[string]string, error) {
// 	folders := make(map[string]string, 0)
// 	files := make(map[string]string, 0)

// 	entries, err := os.ReadDir(path.Join("./uploads", filepath))
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, nil, err
// 	}

// 	for _, e := range entries {
// 		if e.IsDir() {
// 			folders[e.Name()] = path.Join(filepath, e.Name())
// 		} else {
// 			files[e.Name()] = path.Join(filepath, e.Name())
// 		}
// 	}

// 	return folders, files, nil
// }

// func (s *Store) Create(file *multipart.File, filepath string, filename string) error {
// 	newPath := path.Join("uploads", filepath)
// 	err := os.MkdirAll(newPath, os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	dst, err := os.Create(path.Join("uploads", filepath, filename))
// 	if err != nil {
// 		return err
// 	}
// 	defer dst.Close()
// 	if _, err := io.Copy(dst, *file); err != nil {
// 		return err
// 	}
// 	return nil
// }

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
