package to

import (
	"testing"
)

func TestPtr(t *testing.T) {
	b := true
	pb := Ptr(b)
	if pb == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pb != b {
		t.Fatalf("got %v, want %v", *pb, b)
	}
}
