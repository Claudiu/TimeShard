package timeshard

import "testing"

func TestDocument_Save(t *testing.T) {
	c := NewBatch()
	c.Insert(0, []byte("Gandalf"))

	d := &Document{Operations: *c.Snapshot()}

	if err := d.Save("save_test"); err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	d2 := &Document{}
	if err := d2.FromFile("save_test"); err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	if val, _ := d2.Operations.Get(0, MetaDataByteSize); 6 != val {
		t.Fail()
	}
}
