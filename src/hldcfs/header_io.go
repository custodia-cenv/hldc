package hldcfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// SerializeHeader wandelt die Header-Struktur in ein Byte-Array um
func serializeHeader(header *_ImageHeader) ([]byte, error) {
	// Erstellt einen Buffer zum Schreiben der Bytes
	buf := new(bytes.Buffer)

	// Fügt das Präfix "hldc" hinzu
	prefix := []byte("hldc")
	if _, err := buf.Write(prefix); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben des Präfix: %v", err)
	}

	// Schreibt die Version im Little Endian-Format
	if err := binary.Write(buf, binary.LittleEndian, header.Version); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben der Version: %v", err)
	}

	// Schreibt die restlichen Felder im Little Endian-Format
	if err := binary.Write(buf, binary.LittleEndian, header.StartIndexBlock); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben von StartIndexBlock: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, header.IndexSize); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben von IndexSize: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, header.MaxBlocks); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben von MaxBlocks: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, header.MaxSize); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben von MaxSize: %v", err)
	}
	if err := binary.Write(buf, binary.LittleEndian, header.BlockSize); err != nil {
		return nil, fmt.Errorf("fehler beim Schreiben von BlockSize: %v", err)
	}

	// Gibt das Byte-Array zurück
	return buf.Bytes(), nil
}

// DeserializeHeader liest ein Byte-Array und wandelt es in eine Header-Struktur um
func deserializeHeader(data []byte) (*_ImageHeader, error) {
	// Erstellt einen Reader aus dem Byte-Array
	buf := bytes.NewReader(data)

	// Liest und überprüft das Präfix "hldc"
	prefix := make([]byte, 4)
	if _, err := buf.Read(prefix); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen des Präfix: %v", err)
	}
	if string(prefix) != "hldc" {
		return nil, fmt.Errorf("ungültiges Präfix: erwartet 'hldc', erhalten '%s'", string(prefix))
	}

	header := &_ImageHeader{}

	// Liest die Version im Little Endian-Format
	if err := binary.Read(buf, binary.LittleEndian, &header.Version); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen der Version: %v", err)
	}

	// Liest die restlichen Felder im Little Endian-Format
	if err := binary.Read(buf, binary.LittleEndian, &header.StartIndexBlock); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen von StartIndexBlock: %v", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &header.IndexSize); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen von IndexSize: %v", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &header.MaxBlocks); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen von MaxBlocks: %v", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &header.MaxSize); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen von MaxSize: %v", err)
	}
	if err := binary.Read(buf, binary.LittleEndian, &header.BlockSize); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen von BlockSize: %v", err)
	}

	return header, nil
}
