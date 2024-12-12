package rawhldc

import (
	"os"
	"sync"
)

type HldcRawContainer struct {
	mu   *sync.Mutex
	file *os.File
}
