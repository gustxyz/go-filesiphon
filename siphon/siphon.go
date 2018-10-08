package siphon

import (
	"io"
)

//Siphonable type is one that can Get or Put files.
type siphonable interface {
	Get(path string) (io.Reader, error)
	Put(path string, file io.Reader) error
}

//SiphonFile will transfer a file from src Siphonable, at path srcPath to to dest Siphonable at path destPath
func siphonFile(src siphonable, srcPath string, dest siphonable, destPath string) error {
	srcReader, getErr := src.Get(srcPath)

	if getErr != nil {
		return getErr
	}

	err := dest.Put(destPath, srcReader)

	return err
}
