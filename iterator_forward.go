package timeshard

type ForwardIterator struct {
	pos  int64
	over *Block
}

func (it *ForwardIterator) Init() {
	it.pos = -MetaSize
}

func (it *ForwardIterator) Value() []byte {
	if it.GetMeta(MetaOperation) == OpDelete {
		return []byte{}
	}

	meta := it.over.meta[it.pos : it.pos+MetaSize]
	return it.over.data[meta[MetaDataIndex] : meta[MetaDataIndex]+meta[MetaDataByteSize]]
}

func (it *ForwardIterator) Meta() []uint64 {
	return it.over.meta[it.pos : it.pos+MetaSize]
}

func (it *ForwardIterator) GetMeta(metaType MetaType) uint64 {
	meta := it.Meta()
	return meta[metaType]
}

func (it *ForwardIterator) HasNext() bool {
	it.pos += MetaSize
	if int(it.pos) >= len(it.over.meta) {
		return false
	}
	return true
}
