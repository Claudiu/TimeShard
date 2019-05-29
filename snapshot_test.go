package timeshard

import (
	"encoding/json"
	"testing"
	"fmt"
)

func TestSnapshot_MarshalJSON(t *testing.T) {
	c := NewBatch()
	c.Insert(0, []byte("Gandalf"))
	c.Insert(6, []byte(" the "))
	c.Insert(11, []byte("Grey"))
	c.Delete(2, 4)

	b, _ := json.Marshal(c)
	fmt.Println(string(b))
}