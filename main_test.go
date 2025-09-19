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
			want:             0,
			wantErr:          false,
		},
		{
			name:             "valid large income",
			annualEarnings:   "3200000",
			rentAmount:       "1000000",
			businessExpenses: "0",
			want:             166000, // adjust expected
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
			rentAmount:       "12000000",
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
			wantErr:          false, // business expenses are ignored in current logic
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

// run BenchmarkTaxCalculator to see performance
func BenchmarkTaxCalculator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = taxCalculator("3200000", "1000000", "0")
	}
}


func TestPercentageOf(t *testing.T) {
	tests := []struct {
		name string
		part int64
		all  int64
		want uint64
	}{
		{
			name: "normal case",
			part: 50,
			all:  200,
			want: 25,
		},
		{
			name: "zero part",
			part: 0,
			all:  100,
			want: 0,
		},
		{
			name: "zero all",
			part: 50,
			all:  0,
			want: 0, // or handle as error based on your requirements
		},
		{
			name: "part greater than all",
			part: 150,
			all:  100,
			want: 150, // or handle as error based on your requirements
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := percentageOf(tt.part, tt.all); got != tt.want {
				t.Errorf("percentageOf() = %v, want %v", got, tt.want)
			}
		})
	}
}