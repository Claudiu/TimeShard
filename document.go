package timeshard

import (
	"github.com/golang/snappy"
	"github.com/vmihailenco/msgpack"
	"io"
	"os"
	"sync"
)

type Document struct {
	currentSnap *Snapshot

	*sync.RWMutex
}

func (doc *Document) Write(file io.Writer) (int, error) {
	b, err := msgpack.Marshal(doc)
	var out []byte

	if err != nil {
		return 0, err
	}

	out = snappy.Encode(out, b)
	return file.Write(out)
}

func (doc *Document) Save(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	_, err = doc.Write(file)

	return err
}
