package timeshard

type Document struct {
	Title      string            `json:"title"`
	Operations Block             `json:"ops"`
	Meta       map[string]string `json:"meta"`
}

func NewDocument() Document {
	return Document{
		Title:      "Untitled document",
		Operations: *NewBlock(),
		Meta:       map[string]string{},
	}
}

// Bytes returns a slice of length b.Len() holding the end result of our operations
func (doc *Document) Bytes() []byte {
	squashed := doc.Operations.Squash(0)
	for iter := squashed.Iterator(true); iter.HasNext(); {
		return iter.Value()
	}

	return []byte{}
}

// String returns the contents of the document as a string
func (doc *Document) String() string {
	return string(doc.Bytes())
}

func (doc *Document) Insert(at uint64, rawBytes []byte) *Document {
	doc.Operations.Insert(at, rawBytes)
	return doc
}

func (doc *Document) Delete(at uint64, count uint64) *Document {
	doc.Operations.Delete(at, count)
	return doc
}

func (doc *Document) Each(evalFunc func(current Iterator) bool) *Document {
	for iter := doc.Operations.Iterator(false); iter.HasNext(); {
		if evalFunc(iter) != true {
			break
		}
	}

	return doc
}
