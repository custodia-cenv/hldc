package vfs

import (
	"bytes"
	"fmt"
)

// Listet alle Verfügbaren Daten auf
func (o *HldcVfsImage) List() ([]*DataItem, error) {
	return o.index.ListAllDataEntries()
}

// Schreibt Daten unter einem Spezifischen Namen in das VFS
func (o *HldcVfsImage) WriteData(name string, data []byte) error {
	// Die Blöcke werden erzeugt
	dataBlocksk, err := SplitIntoBlocks(data, blockSize)
	if err != nil {
		return err
	}

	// Die Blöcke werden geschrieben
	blockIds := make([]uint64, 0)
	for _, item := range dataBlocksk {
		blockId, err := writeBlockAndGetId(o.raw, item)
		if err != nil {
			return err
		}
		blockIds = append(blockIds, blockId)
	}

	// Es wird eine Kopie des Indexes erzeugt
	indexCopy := o.index.Clone()

	// Es wird ein neuer Eintrag hinzugefügt, dieser Referenziert die Daten
	if err := indexCopy.AddEntry(name, uint64(len(data)), blockIds); err != nil {
		return err
	}

	// Der Neue Index wird geschrieben
	if err := writeNewIndex(o.raw, o.header, indexCopy); err != nil {
		return err
	}

	// Der Index im RAM wird geupdatet
	o.index = indexCopy

	return nil
}

// Gibt Daten zurück
func (o *HldcVfsImage) ReadData(name string) ([]byte, error) {
	// Es wird im Index geprüft ob es Daten unter diesem Namen gibt
	result, foundIt := o.index.entries[name]
	if !foundIt {
		return nil, fmt.Errorf("unkown data")
	}

	// Die Einzelnen Blöcke werden gelesen
	var readBuffer bytes.Buffer
	for _, item := range result.blocks {
		// Der Block wird gelesen
		readedBlock, err := o.raw.ReadBlock(item)
		if err != nil {
			return nil, err
		}

		readBuffer.Write(readedBlock)
	}

	// Die Daten werden extrahiert
	completeData := readBuffer.Bytes()
	readBuffer.Reset()

	// Das Paddenig wird entfernt
	if len(completeData) > int(result.size) {
		completeData = completeData[:result.size]
	}

	// Es wird geprüft ob die Daten vollständig sind
	if len(completeData) != int(result.size) {
		return nil, fmt.Errorf("invalid data")
	}

	return completeData, nil
}

// Gibt die Anazhl der Blöcke innerhalb der des Dateisystemens an
func (o *HldcVfsImage) TotalBlocks() uint64 {
	return o.raw.TotalBlocks()
}

// Schließt das Image Sauber
func (o *HldcVfsImage) Close() error {
	o.raw.Close()
	return nil
}
