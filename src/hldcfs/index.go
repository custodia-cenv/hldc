package hldcfs

func (o *_Index) ListAllDataEntries() ([]*DataItem, error) {
	result := make([]*DataItem, 0)
	for name, item := range o.entries {
		result = append(result, &DataItem{Name: name, Size: item.size})
	}
	return result, nil
}

func (o *_Index) Serialize() ([]byte, error) {
	temp := make([]*_IndexBlockEntry, 0)
	for name, item := range o.entries {
		temp = append(temp, &_IndexBlockEntry{DataName: name, Blocks: item.blocks, StartOffset: 0, TotalSize: item.size})
	}
	bytedIndex, err := serializeIndexBlockEntries(temp...)
	if err != nil {
		return nil, err
	}
	return bytedIndex, nil
}

func (o *_Index) Clone() *_Index {
	newEntries := make(map[string]*_IndexEntry, len(o.entries))
	for key, value := range o.entries {
		newEntries[key] = value // Zeiger werden kopiert, keine neuen _IndexEntry-Objekte
	}
	return &_Index{entries: newEntries}
}

func (o *_Index) AddEntry(name string, size uint64, blcoks []uint64) error {
	o.entries[name] = &_IndexEntry{blocks: blcoks, size: size}
	return nil
}
