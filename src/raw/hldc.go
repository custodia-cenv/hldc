package rawhldc

import (
	"fmt"
	"syscall"
)

func (o *HldcRawContainer) WriteBlock(blockId uint64, data []byte) error {
	// Das Offset wird berechnet
	offset := CalcOffset(blockSize, blockId)

	// Es wird geprüft ob die Daten genauso groß wie die Blockgröße ist sind
	if len(data) != int(blockSize) {
		return fmt.Errorf("data not same size with blocksize: %d != %d", len(data), blockSize)
	}

	// Zum angegebenen Offset springen
	_, err := o.file.Seek(int64(offset), 0)
	if err != nil {
		return fmt.Errorf("fehler beim Setzen des Offsets: %v", err)
	}

	// Daten schreiben
	_, err = o.file.Write(data)
	if err != nil {
		return fmt.Errorf("fehler beim Schreiben der Daten: %v", err)
	}

	return nil
}

func (o *HldcRawContainer) ReadBlock(blockId uint64) ([]byte, error) {
	// Das Offset wird berechnet
	offset := CalcOffset(blockSize, blockId)

	// Zum Start-Offset springen
	_, err := o.file.Seek(offset, 0)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Setzen des Offsets: %v", err)
	}

	// Bytes lesen
	buffer := make([]byte, blockSize)
	n, err := o.file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Lesen der Datei: %v", err)
	}

	// Buffer auf die tatsächliche Anzahl der gelesenen Bytes beschränken
	return buffer[:n], nil
}

func (o *HldcRawContainer) ClearBlock(blockId uint64) error {
	// Erzeugt ein leeres Byte Array
	emptyByteArray := make([]byte, blockSize)

	// Der Leere Block wird geschrieben
	if err := o.WriteBlock(blockId, emptyByteArray); err != nil {
		return err
	}

	return nil
}

func (o *HldcRawContainer) TruncateUpToBlock(blockId uint64) error {
	// Das Offset wird berechnet
	offset := CalcOffset(blockSize, blockId)

	// Kürzen der Datei auf das neue Offset
	err := o.file.Truncate(offset)
	if err != nil {
		return err
	}

	return nil
}

func (o *HldcRawContainer) CopyBlockToAnotherBlock(srcBlockId uint64, destBlockId uint64) error {
	readedBlock, err := o.ReadBlock(srcBlockId)
	if err != nil {
		return err
	}
	if err := o.WriteBlock(destBlockId, readedBlock); err != nil {
		return err
	}
	return nil
}

func (o *HldcRawContainer) TotalBlocks() uint64 {
	// Verwenden Sie die Stat-Methode des *os.File-Objekts
	info, err := o.file.Stat()
	if err != nil {
		panic(err)
	}

	// Dateigröße in Bytes
	size := info.Size()

	// Die Anzahl der Blöcke wird zurückgegben
	return uint64(size / 512)
}

func (o *HldcRawContainer) Close() {
	syscall.Flock(int(o.file.Fd()), syscall.LOCK_UN)
	o.file.Close()
}
