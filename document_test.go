package timeshard

import "testing"

func TestDocument_Save(t *testing.T) {
	d := &Document{}
	if err := d.Save("save_test"); err != nil {
		t.Log(err.Error())
		t.Fail()
	}
}