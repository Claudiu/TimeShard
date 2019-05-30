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
	var out []byte

	if err != nil {
		return 0, err
	}

	out = snappy.Encode(out, b)

	return file.Write(out)
}

func (doc *Document) Save(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	if count, err := file.Write(fileSignature); err != nil || count != len(fileSignature) {
		if count != len(fileSignature) {
			return fmt.Errorf("written bytes differ from signature lenght")
		}

		return err
	}

	_, err = doc.Write(file)

	return err
}

func (doc *Document) FromFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	defer file.Close()

	if err != nil {
		return err
	}

	compressed, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if bytes.Compare(compressed[:len(fileSignature)], fileSignature) != 0 {
		return fmt.Errorf("corupted filed or invalid format")
	}

	return doc.FromBytes(compressed[len(fileSignature):])
}

func (doc *Document) FromBytes(compressed []byte) (err error) {
	var decompressed []byte
	decompressed, err = snappy.Decode(decompressed, compressed)
	if err != nil {
		return err
	}

	temp := &Document{}
	if err := json.Unmarshal(decompressed, &temp); err != nil {
		return err
	}

	doc.Operations = temp.Operations

	return nil
}
