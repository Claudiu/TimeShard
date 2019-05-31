package timeshard

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JSONOperation struct {
	Position uint64 `json:"p"`
	Insert   string `json:"insert,omitempty"`
	Delete   uint64 `json:"delete,omitempty"`
}

// TODO: Unmarshal
func (snap *Snapshot) MarshalJSON() ([]byte, error) {
	iter := snap.Iterator(false)

	var jsonOps []JSONOperation
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

func (snap *Snapshot) UnmarshalJSON(data []byte) error {
	var temp []JSONOperation
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	localBatch := NewBatch()
	for _, op := range temp {
		if op.Insert != "" {
			localBatch.Insert(op.Position, []byte(op.Insert))
			continue
		}

		if op.Delete != 0 {
			localBatch.Delete(op.Position, op.Delete)
			continue
		}


		return fmt.Errorf("could not unmarshal: unknown key")
	}

	copy(snap.meta, localBatch.meta)
	copy(snap.data, localBatch.data)

	return nil
}