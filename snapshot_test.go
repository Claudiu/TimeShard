package timeshard

import (
	"encoding/json"
	"testing"
)

func TestSnapshot_MarshalJSON(t *testing.T) {
	c := NewSnapshot()
	c.Insert(0, []byte("Gandalf"))
	c.Insert(6, []byte(" the "))
	c.Insert(11, []byte("Grey"))
	c.Delete(2, 4)

	_, err := json.Marshal(c)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
