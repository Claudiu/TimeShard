package timeshard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/snappy"
	"io"
	"io/ioutil"
	"os"
)

var fileSignature = []byte{116, 105, 109, 101, 115, 104, 97, 114, 100}

func (doc *Document) Write(file io.Writer) (int, error) {
	b, err := json.Marshal(&doc)

	if err != nil {
		return 0, err
	}

	out := snappy.Encode(nil, b)

	return file.Write(out)
}

func (doc *Document) Save(filename string) (err error) {
	var file *os.File

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// See: https://www.joeshaw.org/dont-defer-close-on-writable-files/
	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	var nBytes int
	nBytes, err = file.Write(fileSignature)
	if err != nil {
		return
	}

	if nBytes != len(fileSignature) {
		return fmt.Errorf("written bytes differ from signature lenght")
	}

	_, err = doc.Write(file)
	if err != nil {
		return
	}

	return
}

func (doc *Document) Open(filename string) (err error) {
	var file *os.File

	file, err = os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		closeError := file.Close()
		if err == nil {
			err = closeError
		}
	}()

	var compressed []byte
	compressed, err = ioutil.ReadAll(file)
	if err != nil {
		return
	}

	sigBytes := compressed[:len(fileSignature)]
	if bytes.Compare(sigBytes, fileSignature) != 0 {
		return fmt.Errorf("corupted filed or invalid format")
	}

	dataBytes := compressed[len(fileSignature):]

	return doc.FromBytes(dataBytes)
}

func (doc *Document) FromBytes(compressed []byte) (err error) {
	decompressed, err := snappy.Decode(nil, compressed)
	if err != nil {
		return err
	}

	temp := &Document{
		Operations: *NewSnapshot(),
	}

	if err := json.Unmarshal(decompressed, &temp); err != nil {
		return err
	}

	doc.Operations = temp.Operations

	return nil
}
