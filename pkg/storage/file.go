package storage

import (
	"log"
	"os"
	"path/filepath"
)

type File interface {
	Store(file string, data []byte) error
}

type file struct {
	dir string
}

func NewFile(dir string) File {
	return &file{dir}
}

// Store implements File.
func (f *file) Store(file string, data []byte) error {
	err := os.MkdirAll(f.dir, 0755)
	if err != nil {
		return err
	}

	fileOut := filepath.Join(f.dir, file)
	err = os.WriteFile(fileOut, data, 0644)
	if err != nil {
		log.Printf("failed to write users to the %q file, err=%q", fileOut, err)
		return err
	}
	return nil
}
