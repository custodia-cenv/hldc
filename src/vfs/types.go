package vfs

import rawhldc "github.com/custodia-cenv/hldc/src/raw"

const blockSize uint16 = 512

type _ImageHeader struct {
	StartIndexBlock uint64 // Gibt die Startposition des Indexes an
	IndexSize       uint64 // Gibt die Gesamtgröße des Indexes an
	MaxBlocks       uint64 // Gibt die Optionale Maximale Anzahl von Blöcken an
	MaxSize         uint64 // Gibt die Optionale Maximale Größe in Bytes an, Größer darf die Datei nicht werden
	BlockSize       uint16 // Gibt die Verwendete Blockgröße an
	Version         uint16 // Gibt die Aktulle Version an
}

type _IndexBlockEntry struct {
	DataName    string   // Speichert den Eigentlichen Dateinamen ab
	Blocks      []uint64 // Speichert die Blöcke Linear ab, in diesen Blöcken befinden sich die eigentlichen Dateien
	StartOffset uint16   // Gibt die Genaue Position innerhalb eines Blocks an, ab wo die Daten eigentlich beginnen
	TotalSize   uint64   // Gibt die Eigentliche Größe der Daten an
}

type _IndexEntry struct {
	blocks []uint64
	size   uint64
}

type _Index struct {
	entries map[string]*_IndexEntry
}

type HldcVfsImage struct {
	header *_ImageHeader
	index  *_Index
	raw    *rawhldc.HldcRawContainer
}

type DataItem struct {
	Name string
	Size uint64
}

type VfsFile struct {
}

type VfsFileInfo struct {
}
