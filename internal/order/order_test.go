package order

import (
	"errors"
	"math"
	"testing"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		want    float64
		wantErr error
	}{
		{
			name: "gold customer gets 20% off",
			order: &Order{ID: 1, Customer: &Customer{Tier: "gold"},
				Items: []Item{{Price: 100, Quantity: 1}}},
			want: 80,
		},
		{
			name: "no tier pays full price",
			order: &Order{ID: 2, Customer: &Customer{},
				Items: []Item{{Price: 10, Quantity: 3}}},
			want: 30,
		},
		{
			name:    "empty order is an error",
			order:   &Order{ID: 3, Customer: &Customer{}},
			wantErr: ErrEmptyOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Process(tt.order)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("total = %v, want %v", got, tt.want)
			}
		})
	}
}
