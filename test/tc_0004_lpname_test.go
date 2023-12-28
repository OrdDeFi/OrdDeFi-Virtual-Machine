package test

import (
	"OrdDeFi-Virtual-Machine/virtual_machine/memory/memory_write"
	"testing"
)

func TestLPName(t *testing.T) {
	r := memory_write.ODFISpendingTickLPName("odgv")
	if *r != "odfi-odgv" {
		t.Errorf("calc LP name error, got %s, expected %s", *r, "odfi-odgv")
	}
	r = memory_write.ODFISpendingTickLPName("aaaa")
	if *r != "aaaa-odfi" {
		t.Errorf("calc LP name error, got %s, expected %s", *r, "aaaa-odfi")
	}
	r = memory_write.ODFISpendingTickLPName("odfi")
	if r != nil {
		t.Errorf("calc LP name error, got %s, expected nil", *r)
	}
}
