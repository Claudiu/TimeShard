package timeshard

import "sync"

type Document struct {
	currentSnap *Snapshot

	*sync.RWMutex
}

//func (document *Document) PushBatch(batch *Batch) {
//	document.Lock()
//	document.Batches = append(document.Batches, batch)
//	document.Unlock()
//}

