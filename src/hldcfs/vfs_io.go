package hldcfs

import (
	"fmt"
	"io/fs"
	"sync"
)

func Open(image *HldcVfsImage, name string) (*VfsFile, error) {
	// Der Mutex wird verwendet
	image.mu.Lock()
	defer image.mu.Unlock()

	// Es wird im Index geprüft ob es Daten unter diesem Namen gibt
	result, foundIt := image.index.entries[name]
	if !foundIt {
		return nil, fmt.Errorf("unkown data")
	}

	// Das Rückgabe Image wird erzeugt
	returnImage := &VfsFile{
		mu:    new(sync.Mutex),
		size:  result.size,
		image: image, blocks: result.blocks,
		name: name,
	}

	// Das Image wird zurückgegben
	return returnImage, nil
}

func Create(image *HldcVfsImage, name string) (*VfsFile, error) {
	// Der Mutex wird verwendet
	image.mu.Lock()
	defer image.mu.Unlock()

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
