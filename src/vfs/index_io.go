package vfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func serializeIndexBlockEntries(entries ...*_IndexBlockEntry) ([]byte, error) {
	var buf bytes.Buffer

	// Schreibe die Anzahl der Einträge als uint64
	numEntries := uint64(len(entries))
	if err := binary.Write(&buf, binary.LittleEndian, numEntries); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben der Anzahl der Einträge: %v", err)
	}

	// Iteriere über jedes IndexEntry und schreibe die Daten
	for _, entry := range entries {
		// Serialisiere DataName
		fileNameBytes := []byte(entry.DataName)
		fileNameLength := uint16(len(fileNameBytes))
		if err := binary.Write(&buf, binary.LittleEndian, fileNameLength); err != nil {
			return nil, fmt.Errorf("fehler beim Schreiben der DataName-Länge: %v", err)
		}
		if _, err := buf.Write(fileNameBytes); err != nil {
			return nil, fmt.Errorf("fehler beim Schreiben der DataName: %v", err)
		}

		// Serialisiere Blocks
		numBlocks := uint32(len(entry.Blocks))
		if err := binary.Write(&buf, binary.LittleEndian, numBlocks); err != nil {
			return nil, fmt.Errorf("fehler beim Schreiben der Anzahl der Blocks: %v", err)
		}
		for _, block := range entry.Blocks {
			if err := binary.Write(&buf, binary.LittleEndian, block); err != nil {
				return nil, fmt.Errorf("fehler beim Schreiben eines Blocks: %v", err)
			}
		}

		// Serialisiere StartOffset
		if err := binary.Write(&buf, binary.LittleEndian, entry.StartOffset); err != nil {
			return nil, fmt.Errorf("fehler beim Schreiben des StartOffset: %v", err)
		}

		// Serialisiere TotalSize
		if err := binary.Write(&buf, binary.LittleEndian, entry.TotalSize); err != nil {
			return nil, fmt.Errorf("fehler beim Schreiben des TotalSize: %v", err)
		}
	}

	return buf.Bytes(), nil
}

func deserializeIndexBlockEntries(data []byte) ([]*_IndexBlockEntry, error) {
	var entries []*_IndexBlockEntry
	reader := bytes.NewReader(data)

	// Lese die Anzahl der Einträge (uint64)
	var numEntries uint64
	if err := binary.Read(reader, binary.LittleEndian, &numEntries); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen der Anzahl der Einträge: %v", err)
	}

	// Iteriere über die Anzahl der Einträge und lese jeden Eintrag
	for i := uint64(0); i < numEntries; i++ {
		entry := &_IndexBlockEntry{}

		// Lese die Länge von DataName (uint16)
		var fileNameLength uint16
		if err := binary.Read(reader, binary.LittleEndian, &fileNameLength); err != nil {
			return nil, fmt.Errorf("fehler beim Lesen der DataName-Länge: %v", err)
		}

		// Lese die DataName Bytes
		fileNameBytes := make([]byte, fileNameLength)
		if _, err := io.ReadFull(reader, fileNameBytes); err != nil {
			return nil, fmt.Errorf("fehler beim Lesen der DataName: %v", err)
		}
		entry.DataName = string(fileNameBytes)

		// Lese die Anzahl der Blocks (uint32)
		var numBlocks uint32
		if err := binary.Read(reader, binary.LittleEndian, &numBlocks); err != nil {
			return nil, fmt.Errorf("fehler beim Lesen der Anzahl der Blocks: %v", err)
		}

		// Lese die Blocks als []uint64
		entry.Blocks = make([]uint64, numBlocks)
		for j := uint32(0); j < numBlocks; j++ {
			var block uint64
			if err := binary.Read(reader, binary.LittleEndian, &block); err != nil {
				return nil, fmt.Errorf("fehler beim Lesen des Blocks %d: %v", j, err)
			}
			entry.Blocks[j] = block
		}

		// Lese StartOffset (uint64)
		if err := binary.Read(reader, binary.LittleEndian, &entry.StartOffset); err != nil {
			return nil, fmt.Errorf("fehler beim Lesen des StartOffset: %v", err)
		}

		// Lese TotalSize (uint64)
		if err := binary.Read(reader, binary.LittleEndian, &entry.TotalSize); err != nil {
			return nil, fmt.Errorf("fehler beim Lesen des TotalSize: %v", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
