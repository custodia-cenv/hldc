package rawhldc

func CalcOffset(blockSize uint16, blockId uint64) int64 {
	// Der Startoffset wird anahnd der BlockID errechnet
	var stOffset int64
	if blockId == 0 {
		stOffset = 0
	} else if blockId == 1 {
		stOffset = int64(blockSize)
	} else {
		stOffset = int64(blockSize) * int64(blockId)
	}
	return stOffset
}
