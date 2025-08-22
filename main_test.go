package main

import (
	"testing"
)

func TestTaxCalculator(t *testing.T) {
	tests := []struct {
		name             string
		annualEarnings   string
		rentAmount       string
		businessExpenses string
		want             uint64
		wantErr          bool
	}{
		{
			name:             "valid small income",
			annualEarnings:   "500000",
			rentAmount:       "100000",
			businessExpenses: "50000",
			want:             350000, // adjust based on your formula
			wantErr:          false,
		},
		{
			name:             "valid large income",
			annualEarnings:   "2200000",
			rentAmount:       "300000",
			businessExpenses: "150000",
			want:             1750000, // adjust expected
			wantErr:          false,
		},
		{
			name:             "zero income",
			annualEarnings:   "0",
			rentAmount:       "0",
			businessExpenses: "0",
			want:             0,
			wantErr:          false,
		},
		{
			name:             "invalid annual earnings",
			annualEarnings:   "abc",
			rentAmount:       "20000",
			businessExpenses: "10000",
			want:             0,
			wantErr:          true,
		},
		{
			name:             "invalid rent",
			annualEarnings:   "500000",
			rentAmount:       "xyz",
			businessExpenses: "10000",
			want:             0,
			wantErr:          true,
		},
		{
			name:             "invalid business expenses",
			annualEarnings:   "500000",
			rentAmount:       "20000",
			businessExpenses: "oops",
			want:             0,
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := taxCalculator(tt.annualEarnings, tt.rentAmount, tt.businessExpenses)
			if (err != nil) != tt.wantErr {
				t.Errorf("taxCalculator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("taxCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}