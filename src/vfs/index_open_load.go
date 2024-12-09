package vfs

func newIndex() *_Index {
	return &_Index{
		entries: make(map[string]*_IndexEntry),
	}
}

func loadFromBytes(data []byte) (*_Index, error) {
	// Der Index wird eingelesen
	indexEntries, err := deserializeIndexBlockEntries(data)
	if err != nil {
		return nil, err
	}

	// Es wird ein neuer Leerer Index erezugt
	emptyIndex := newIndex()

	// Die Einträge werden hinzugefügt
	for _, item := range indexEntries {
		emptyIndex.entries[item.DataName] = &_IndexEntry{blocks: item.Blocks, size: item.TotalSize}
	}

	// Der Neue Index wird zurückgegeben
	return emptyIndex, nil
}
