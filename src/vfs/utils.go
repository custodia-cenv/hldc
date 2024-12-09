package vfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func SplitIntoBlocks(data []byte, blockSize uint16) ([][]byte, error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("blockSize muss größer als 0 sein")
	}

	var returnSlice [][]byte

	for len(data) > int(blockSize) {
		// Füge einen vollständigen Block hinzu
		returnSlice = append(returnSlice, data[:blockSize])
		// Schneide die verarbeiteten Daten ab
		data = data[blockSize:]
	}

	// Verarbeite den letzten Block, der kleiner oder gleich blockSize ist
	if len(data) > 0 {
		// Erstelle ein Padding, wenn der letzte Block kleiner als blockSize ist
		padding := make([]byte, int(blockSize)-len(data))
		paddedData := append(data, padding...)
		returnSlice = append(returnSlice, paddedData)
	}

	return returnSlice, nil
}

func SplitIntoWithNextLink(data []byte, blockSize uint16, currentBlock uint64) ([][]byte, error) {
	// Überprüfen, ob blockSize groß genug ist, um den Header aufzunehmen
	if blockSize <= 8 {
		return nil, fmt.Errorf("blockSize muss größer als 8 sein")
	}

	var returnSlice [][]byte
	offset := 0
	nextBlockId := currentBlock + 1

	// Anzahl der Datenbytes pro Block (blockSize minus 8 Bytes für den Header)
	dataPerBlock := blockSize - 8

	for offset < len(data) {
		buf := new(bytes.Buffer)

		// Bestimmen, ob weitere Blöcke folgen
		isLastBlock := (offset + int(dataPerBlock)) >= len(data)

		// Daten für den aktuellen Block extrahieren
		end := offset + int(dataPerBlock)
		if end > len(data) {
			end = len(data)
		}
		chunk := data[offset:end]

		// Wenn die Daten weniger als dataPerBlock Bytes enthalten, mit Nullen auffüllen
		if len(chunk) < int(dataPerBlock) {
			padding := make([]byte, int(dataPerBlock)-len(chunk))
			chunk = append(chunk, padding...)
		}

		// Daten in den Buffer schreiben
		if _, err := buf.Write(chunk); err != nil {
			return nil, fmt.Errorf("fehler beim Schreiben der Daten: %v", err)
		}

		// Die ID des nächsten Blocks schreiben, sofern es sich nicht um den letzten handelt
		if !isLastBlock {
			if err := binary.Write(buf, binary.LittleEndian, nextBlockId); err != nil {
				return nil, fmt.Errorf("fehler beim Schreiben der Block-ID: %v", err)
			}
		} else {
			// Acht Nullen für den letzten Block schreiben
			padding := make([]byte, 8)
			if _, err := buf.Write(padding); err != nil {
				return nil, fmt.Errorf("fehler beim Schreiben der Padding-Bytes: %v", err)
			}
		}

		// Den fertigen Block zum Rückgabeslice hinzufügen
		returnSlice = append(returnSlice, buf.Bytes())

		// Offset und Block-ID für die nächste Iteration aktualisieren
		offset += int(dataPerBlock)
		nextBlockId++
	}

	return returnSlice, nil
}
