package hldcfs

import (
	"fmt"

	rawhldc "github.com/custodia-cenv/hldc/src/raw"
)

func writeHedader(raw *rawhldc.HldcRawContainer, header *_ImageHeader) error {
	// Der Header wird in Bytes umgewaldet
	bytesHeader, err := serializeHeader(header)
	if err != nil {
		return err
	}

	// Der Header wird in Blöcke umgewandelt
	bytesHeaderBlocks, err := SplitIntoBlocks(bytesHeader, blockSize)
	if err != nil {
		return err
	}

	// Es muss Explizit ein 1 Block sein
	if len(bytesHeaderBlocks) != 1 {
		return fmt.Errorf("invalid header")
	}

	// Der Header Block wird geschrieben
	if err := raw.WriteBlockHDD(0, bytesHeaderBlocks[0]); err != nil {
		return err
	}

	return nil
}

func updateHeader(raw *rawhldc.HldcRawContainer, cheader *_ImageHeader, indexSize uint64, indexStartBlock uint64) error {
	// Der Header wird geupdated
	updatedHeader := &_ImageHeader{StartIndexBlock: indexStartBlock, IndexSize: indexSize, MaxBlocks: cheader.MaxBlocks, MaxSize: cheader.MaxSize, Version: cheader.Version}

	// Der Header wird geschrieben
	if err := writeHedader(raw, updatedHeader); err != nil {
		return err
	}
	return nil
}

func writeNewIndex(raw *rawhldc.HldcRawContainer, cheader *_ImageHeader, newIndex *_Index) error {
	// Der Index wird in Bytes umgewandelt
	bytedIndex, err := newIndex.Serialize()
	if err != nil {
		return err
	}

	// Die Anzahl der Blöcke wird ermittelt
	startHight := raw.TotalBlocks()
	blockHight := startHight

	// Die Bytes werden in Blöcke aufgeteilt
	splitedBlocksWithNextBlockLink, err := SplitIntoWithNextLink(bytedIndex, blockSize, blockHight)
	if err != nil {
		return err
	}

	// Die Blöcke werden nach und nach geschrieben
	for _, item := range splitedBlocksWithNextBlockLink {
		if err := raw.WriteBlockHDD(blockHight, item); err != nil {
			return err
		}
		blockHight++
	}

	// Der Header wird geupdated
	if err := updateHeader(raw, cheader, uint64(len(bytedIndex)), startHight); err != nil {
		return err
	}

	return nil
}

func writeBlockAndGetId(raw *rawhldc.HldcRawContainer, data []byte) (uint64, error) {
	// Die Aktuelle Blockhöhe wird ermittelt
	currentBlockHight := raw.TotalBlocks()

	// Der Block wird geschrieben
	if err := raw.WriteBlockHDD(currentBlockHight, data); err != nil {
		return 0, err
	}

	// Die ID wird zurückgegben
	return currentBlockHight, nil
}
