package vfs

import (
	"io/fs"
)

func Open(image *HldcVfsImage, name string) (*VfsFile, error) {
	return &VfsFile{}, nil
}

func Create(image *HldcVfsImage, name string) (*VfsFile, error) {
	return nil, nil
}

func Remove(image *HldcVfsImage, name string) error {
	return nil
}

func Rename(image *HldcVfsImage, oldName string, NewName string) error {
	return nil
}

func Stat(image *HldcVfsImage, name string) (*VfsFileInfo, error) {
	return nil, nil
}

func Truncate(image *HldcVfsImage, name string, size uint64) error {
	return nil
}

func OpenFile(image *HldcVfsImage, name string, flag int, perm fs.FileMode) (*VfsFile, error) {
	return nil, nil
}

func ReadFile(image *HldcVfsImage, name string) ([]byte, error) {
	return nil, nil
}

func WriteFile(image *HldcVfsImage, name string, data []byte, perm fs.FileMode) error {
	return nil
}
