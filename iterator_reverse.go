package timeshard

type ReverseIterator struct {
	pos  int64
	over *Snapshot
}

func (it *ReverseIterator) Init() {
	it.pos = int64(len(it.over.meta) + MetaSize)
}

func (it *ReverseIterator) Value() []byte {
	if it.GetMeta(MetaOperation) == OpDelete {
		return []byte{}
	}

	meta := it.over.meta[it.pos-MetaSize : it.pos]
	return it.over.data[meta[MetaDataIndex] : meta[MetaDataIndex]+meta[MetaDataByteSize]]
}

func (it *ReverseIterator) Meta() []uint64 {
	return it.over.meta[it.pos-MetaSize : it.pos]
}

func (it *ReverseIterator) GetMeta(metaType MetaType) uint64 {
	meta := it.Meta()
	return meta[metaType]
}

func (it *ReverseIterator) AffectedArea() uint64 {
	meta := it.over.meta[it.pos-MetaSize : it.pos]
	return meta[MetaDataByteSize]
}

func (it *ReverseIterator) PointInTime() uint64 {
	meta := it.over.meta[it.pos-MetaSize : it.pos]
	return meta[MetaTimestamp]
}

func (it *ReverseIterator) Type() uint64 {
	meta := it.over.meta[it.pos-MetaSize : it.pos]
	return meta[MetaOperation]
}

func (it *ReverseIterator) Retain() uint64 {
	meta := it.over.meta[it.pos-MetaSize : it.pos]
	return meta[MetaRetain]
}

func (it *ReverseIterator) HasNext() bool {
	it.pos = it.pos - MetaSize

	if int(it.pos) <= 0 {
		return false
	}

	return true
}
