package timeshard

import (
	"bytes"
	"errors"
	"sync"
	"time"
	"unicode/utf8"
)

type Block struct {
	Shard
	sync.Mutex
}

func NewBlock() *Block {
	return &Block{
		Shard: Shard{
			make([]byte, 0),
			make([]uint64, 0),
		},
		Mutex: sync.Mutex{},
	}
}

func (block *Block) Iterator(reverse bool) Iterator {
	var iterator Iterator

	if reverse {
		iterator = &ReverseIterator{
			over: block.Clone(),
		}

		iterator.Init()
		return iterator
	}

	iterator = &ForwardIterator{
		over: block.Clone(),
	}

	iterator.Init()
	return iterator
}

func (block *Block) LastActivity() uint64 {
	data, _ := block.Get(0, MetaTimestamp)
	return data
}

func (block *Block) Squash(count uint64) *Block {
	//newSnap := NewEmpty()

	mirror := ""

	iter := block.Iterator(false)

	current := uint64(0)
	for iter.HasNext() && (current < count || count == 0) {
		maxLen := uint64(utf8.RuneCountInString(mirror))

		switch iter.GetMeta(MetaOperation) {
		case OpInsert:
			pos := iter.GetMeta(MetaRetain)
			if pos > maxLen {
				pos = maxLen
			}

			n := string(bytes.Runes(iter.Value()))
			mirror = mirror[:pos] + n + mirror[pos:]
		case OpDelete:
			pos := iter.GetMeta(MetaRetain)
			if pos > maxLen {
				continue
			}

			affected := iter.GetMeta(MetaDataByteSize)
			if pos+affected > maxLen {
				affected = maxLen - pos
			}

			mirror = mirror[:pos] + mirror[pos+affected:]
		}

		current++
	}

	b := NewBlock()
	b.Insert(0, []byte(mirror))

	return b
}

// add will add an operation into our Shard
func (block *Block) add(rawBytes []byte, action, retain uint64) {
	block.Lock()
	defer block.Unlock()

	var s, l uint64
	if foundIndex := bytes.Index(block.data, rawBytes); foundIndex != -1 {
		s = uint64(foundIndex)
		l = uint64(len(rawBytes))
	} else {
		s, l = block.pushData(&rawBytes)
	}

	block.pushMeta(s, l, retain, OpInsert)
}

// Insert will add an OpInsert into our Shard
func (block *Block) Insert(at uint64, rawBytes []byte) {
	block.add(rawBytes, OpInsert, at)
}

// Delete will add an OpDelete into our Shard
// TODO: Use add
func (block *Block) Delete(at uint64, count uint64) {
	block.Lock()
	defer block.Unlock()

	now := time.Now().UnixNano()

	currentLength := uint64(len(block.data))
	data := []uint64{currentLength, count, OpDelete, at, uint64(now)}
	block.meta = append(block.meta, data...)
}

// Clone will clone the current snapshot, useful for iterating
func (block *Block) Clone() *Block {
	block.Lock()
	defer block.Unlock()

	localBlock := NewBlock()

	localBlock.Shard = Shard{
		make([]byte, len(block.data)),
		make([]uint64, len(block.meta)),
	}

	copy(localBlock.data, block.data)
	copy(localBlock.meta, block.meta)

	return localBlock
}

// Will merge the instructions from two Snapshots
func Merge(src *Block, dst *Block) (*Block, error) {
	if integrity := src.assertIntegrity(); integrity {
		return nil, errors.New("first Batch failed integrity check")
	}

	if integrity := dst.assertIntegrity(); integrity {
		return nil, errors.New("second Batch failed integrity check")
	}

	// Verify whether src is older then dst
	first, last := &src, &dst
	if src.LastActivity() > dst.LastActivity() {
		first, last = &dst, &src
	}

	// Data does not need to be moved...
	//for i := 0; i < len((*last).meta); i += MetaSize {
	//	(*last).meta[i] += uint64(len((*first).data))
	//}

	dataMerged := append((*first).data, (*last).data...)
	metaMerged := append((*first).meta, (*last).meta...)

	var targetSnapshot *Block

	if dst == nil {
		targetSnapshot = NewBlock()
	} else {
		targetSnapshot = dst
	}

	targetSnapshot.data = append(targetSnapshot.data, dataMerged...)
	targetSnapshot.meta = append(targetSnapshot.meta, metaMerged...)

	return targetSnapshot, nil
}
