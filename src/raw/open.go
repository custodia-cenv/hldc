package rawhldc

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

func CreateNewHLDCDataContainerRAW(filename string) error {
	// Datei erstellen (oder überschreiben)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func CreateNewHLDCDataContainerAndOpenRAW(filename string) (*HldcRawContainer, error) {
	if err := CreateNewHLDCDataContainerRAW(filename); err != nil {
		return nil, err
	}
	HldcRawContainer, err := OpenHLDCDataContainerRAW(filename)
	if err != nil {
		return nil, err
	}
	return HldcRawContainer, nil
}

func OpenHLDCDataContainerRAW(filename string) (*HldcRawContainer, error) {
	// Datei im Lesemodus öffnen
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	// Sperren Sie die Datei
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return nil, err
	}

	// Verwenden Sie die Stat-Methode des *os.File-Objekts
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// Dateigröße in Bytes
	size := info.Size()

	// Es ist kein Index vorhanden, es handelt sich um ein Leeres Image
	rvobj := &HldcRawContainer{
		mu:   new(sync.Mutex),
		file: file,
	}

	// Überprüfen, ob die Größe ein Vielfaches von 512 ist
	if uint64(size)%uint64(blockSize) != 0 {
		fmt.Printf("Die Datei verstößt gegen die 512-Byte-Regel: %d ist nicht durch %d teilbar.\n", size, blockSize)
	}

	// Das Objekt wird zurückgegen
	return rvobj, nil
}
