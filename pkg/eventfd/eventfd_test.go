package eventfd

import (
	"testing"
)

func TestNew(t *testing.T) {
	efd, err := New()
	if err != nil {
		t.Error("Could not create EventFD")
	}
	defer efd.Close()

	if efd.Fd() == 0 {
		t.Errorf("Invalid FD %d", efd.Fd())
	}

}

func TestReadWrite(t *testing.T) {
	efd, err := New()
	if err != nil {
		t.Errorf("%q", err)
	}
	defer efd.Close()

	// write val
	buf := make([]byte, 8)
	buf[0] = 0x01
	n, err := efd.Write(buf)
	if err != nil {
		t.Errorf("Could not write to eventfd %q", err)
	}
	if n != 8 {
		t.Error("Error while writing to eventfd")
	}

	buf[0] = 0
	n, err = efd.Read(buf)
	if err != nil {
		t.Errorf("Could not read from eventfd %q", err)
	}
	if n != 8 {
		t.Error("Error while reading from eventfd")
	}
	if buf[0] != 0x01 {
		t.Errorf("Expected 0x01 found %x", buf[0])
	}
}

func TestReadWriteEvents(t *testing.T) {
	efd, err := New()
	if err != nil {
		t.Error(err)
	}
	defer efd.Close()
	var val, good uint64
	good = 0x0011223344556677

	val = good
	err = efd.WriteEvents(val)
	if err != nil {
		t.Error(err)
	}

	val, err = efd.ReadEvents()
	if err != nil {
		t.Error(err)
	}
	if val != good {
		t.Errorf("Error while reading from eventfd, expected %q got %q", good, val)
	}
}
