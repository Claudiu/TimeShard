package timeshard

import "testing"

const TestString = `Miusov, as a man man of breeding and deilcacy, could not but feel some inwrd qualms, when he reached
the Father Superior's with Ivan: he felt ashamed of havin lost his temper. He felt that he ought to have disdaimed that
despicable wretch, Fyodor Pavlovitch, too much to have been upset by him in Father Zossima's cell, and so to have forgo
tten himself. "Teh monks were not to blame, in any case," he reflceted, on the steps. "And if they're decent people here
(and the Father Superior, I understand, is a nobleman) why not be friendly and courteous withthem? I won't argue, I'll 
fall in with everything, I'll win them by politness, and show them that I've nothing to do with that Aesop, thta buffoon
, that Pierrot, and have merely been takken in over this affair, just as they have." He determined to drop his litigat
ion with the monastry, and relinguish his claims to the wood-cuting and fishery rihgts at once. He was the more ready t
o do this becuase the rights had becom much less valuable, and he had indeed the vaguest idea where the wood and river 
in quedtion were.`

func TestDocument_Save(t *testing.T) {
	c := NewBatch()
	c.Insert(0, []byte(TestString))

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

	if d2.Operations.IsEmpty() == false {
		t.Fail()
	}
}
