package timeshard

import (
	"encoding/json"
	"testing"
)

func TestSnapshot_MarshalJSON(t *testing.T) {
	c := NewBlock()
	c.Insert(0, []byte("Gandalf"), nil)
	c.Insert(6, []byte(" the "), nil)
	c.Insert(11, []byte("Grey"), nil)
	c.Delete(2, 4)

	_, err := json.Marshal(c)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
