package main

import (
	"testing"
)


func TestTaxCalculator(t *testing.T) {

	tax, err := taxCalculator("700000.00", "0.00", "0.00")
	if err != nil {
		t.Error(err)
	}

	expected := uint64(0)
	if expected != tax {
		t.Errorf("expected: %d; got: %d;", expected, tax)
	}
	
}