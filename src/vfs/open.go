package vfs

import (
	"bytes"
	"encoding/binary"
	"fmt"

	rawhldc "github.com/custodia-cenv/hldc/src/raw"
)

func CreateNewHldcVfsImage(filename string, maxBlocks uint64, maxSize uint64) error {
	// Es wird eine neue RAW Datei erzeugt
	blockImage, err := rawhldc.CreateNewHLDCDataContainerAndOpenRAW(filename)
	if err != nil {
		return err
	}

	// Der Header wird erzeugt
	newHeader := &_ImageHeader{
		StartIndexBlock: 0,
		IndexSize:       0,
		MaxBlocks:       maxBlocks,
		MaxSize:         maxSize,
		Version:         1,
	}

	// Der Header wird geschrieben
	if err := writeHedader(blockImage, newHeader); err != nil {
		return err
	}

	// Es muss sich genau 1 Block in der Datei befinden
	if blockImage.TotalBlocks() != 1 {
		return fmt.Errorf("invalid image created, unkown reason")
	}

	// Die Datei wird geschlossen
	blockImage.Close()

	// Der Vorgang war erfolgreich
	return nil
}

func CreateNewHldcVfsImageAndOpen(filename string, maxBlocks uint64, maxSize uint64) (*HldcVfsImage, error) {
	if err := CreateNewHldcVfsImage(filename, maxBlocks, maxSize); err != nil {
		return nil, err
	}
	HldcRawContainer, err := OpenHldcVfsImage(filename)
	if err != nil {
		return nil, err
	}
	return HldcRawContainer, nil
}

func OpenHldcVfsImage(filename string) (*HldcVfsImage, error) {
	// Die RAW Datei wird geöffnet
	rawImage, err := rawhldc.OpenHLDCDataContainerRAW(filename)
	if err != nil {
		return nil, err
	}

	// Es wird geprüft ob mindestens 1 Block vorhanden ist
	if rawImage.TotalBlocks() < 1 {
		return nil, fmt.Errorf("invalid hldc-vfs image")
	}

	// Block 0 wird geladen
	headerBlock, err := rawImage.ReadBlock(0)
	if err != nil {
		return nil, err
	}

	// Es wird versucht den Header einzulesen
	header, err := deserializeHeader(headerBlock[:40])
	if err != nil {
		return nil, err
	}

	// Es wird geprüft ob es bereits einen Index gibt
	var imageIndex *_Index
	if header.StartIndexBlock == 0 {
		if header.IndexSize != 0 {
			return nil, fmt.Errorf("broken image")
		}
		imageIndex = newIndex()
	} else {
		// Es werden alle Index Blöcke abgerufen
		var readBuffer bytes.Buffer
		nextReadingBlock := header.StartIndexBlock
		for readBuffer.Len() != int(header.IndexSize) {
			// Sollte der nächste Block 0 sein, wird abgebrochen
			if nextReadingBlock == 0 {
				break
			}

			// Der Indexblock wird gelesen
			readedBlock, err := rawImage.ReadBlock(nextReadingBlock)
			if err != nil {
				return nil, err
			}

			// Die letzten 8 Bytes werden abgeschnitten
			lastEightBytes := readedBlock[len(readedBlock)-8:]
			nextReadingBlock = binary.LittleEndian.Uint64(lastEightBytes)

			// Der Block ohne die letzten 8 Bytes wird extrahiert
			readedBlock = readedBlock[:len(readedBlock)-8]

			if len(readedBlock) > int(header.IndexSize) {
				readBuffer.Write(readedBlock[:header.IndexSize])
			} else if (len(readedBlock) + readBuffer.Len()) > int(header.IndexSize) {
				sSize := uint(header.IndexSize) - uint(len(readedBlock)+readBuffer.Len())
				readBuffer.Write(readedBlock[:sSize])
			} else {
				readBuffer.Write(readedBlock)
			}
		}

		// Der Index wird wiederhergestellt
		imageIndex, err = loadFromBytes(readBuffer.Bytes())
		if err != nil {
			return nil, err
		}
	}

	// Es wird ein HldcVfsImage Objekt erzeugt
	rvobj := &HldcVfsImage{
		header: header,
		index:  imageIndex,
		raw:    rawImage,
	}

	// Das Objekt wird zurückgegen
	return rvobj, nil
}
