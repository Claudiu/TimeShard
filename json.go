package timeshard

import (
	"bytes"
	"encoding/json"
)

type JSONOperation struct {
	Position uint64 `json:"p"`
	Insert   string `json:"insert,omitempty"`
	Delete   uint64 `json:"delete,omitempty"`
}

// TODO: Unmarshal
func (snap *Snapshot) MarshalJSON() ([]byte, error) {
	iter := snap.Iterator(false)

	var jsonOps = []JSONOperation{}
	for iter.HasNext() {
		switch iter.GetMeta(MetaOperation) {
		case OpInsert:
			str := string(bytes.Runes(iter.Value()))
			jsonOps = append(jsonOps, JSONOperation{
				Position: iter.GetMeta(MetaRetain),
				Insert:   str,
			})

		case OpDelete:
			jsonOps = append(jsonOps, JSONOperation{
				Position: iter.GetMeta(MetaRetain),
				Delete:   iter.GetMeta(MetaDataByteSize),
			})
		}
	}

	return json.Marshal(jsonOps)
}
