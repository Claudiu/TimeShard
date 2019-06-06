package timeshard

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/golang/snappy"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
)

var fileSignature = []byte("application/timeshard")
var crc32Table = crc32.MakeTable(crc32.Koopman)

const crc32Bytes = 8

func computeCRC32(b []byte) []byte {
	crcSum := crc32.Checksum(b, crc32Table)

	crcBytes := make([]byte, crc32Bytes)
	binary.LittleEndian.PutUint32(crcBytes, crcSum)

	return crcBytes
}

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
		return errors.New("written bytes differ from signature length")
	}

	var buf bytes.Buffer
	_, err = doc.Write(&buf)
	if err != nil {
		return
	}

	// we compute a CRC32 Sum using CRC32.Koopman
	_, err = file.Write(computeCRC32(buf.Bytes()))
	if err != nil {
		return
	}

	_, err = file.Write(buf.Bytes())
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
		return errors.New("corrupted filed or invalid format")
	}

	dataBytes := compressed[crc32Bytes+len(fileSignature):]

	crc32FromFile := compressed[len(fileSignature) : len(fileSignature)+crc32Bytes]
	if bytes.Compare(crc32FromFile, computeCRC32(dataBytes)) != 0 {
		return errors.New("CRC32 hash mismatch")
	}

	return doc.FromBytes(dataBytes)
}

func (doc *Document) FromBytes(compressed []byte) (err error) {
	decompressed, err := snappy.Decode(nil, compressed)
	if err != nil {
		return err
	}

	temp := &Document{
		Operations: *NewBlock(),
	}

	if err := json.Unmarshal(decompressed, &temp); err != nil {
		return err
	}

	doc.Operations = temp.Operations

	return nil
}
