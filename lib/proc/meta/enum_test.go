package meta

import "testing"

func TestEnumDefaultValues(t *testing.T) {
	// important as structs will default integers to 0
	if OptBoolNone != 0 {
		t.Error("OptBoolNone will not be default!")
	}
	if ProtoNone != 0 {
		t.Error("ProtoNone will not be default!")
	}
}
