package mr

import "testing"

func TestA(t *testing.T) {

	n := newFileStatus([]string{"1", "2", "3"})
	a, b := n.getJob()
	if b {
		t.Errorf("b")
	}
	c, b := n.getJob()
	if b {
		t.Errorf("b")
	}
	n.addWork(a, 1)
	n.addWork(c, 2)
	n.workDone(1)
	n.workDone(2)
	_, job := n.getJob()
	if !job {
		t.Errorf("should no work")
	}
}
