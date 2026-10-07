// Package order contains a tiny order-processing pipeline used to
// demonstrate debugging and stack tracing in Go.
package order

import (
	"errors"
	"fmt"
)

// Item is a single line in an order.
type Item struct {
	SKU      string
	Price    float64
	Quantity int
}

// Customer places orders. Tier controls the discount they receive.
type Customer struct {
	Name string
	Tier string // "gold", "silver" or "" (none)
}

// Order groups a customer with the items they bought.
type Order struct {
	ID       int
	Customer *Customer
	Items    []Item
}

// ErrEmptyOrder is returned when an order has no items.
var ErrEmptyOrder = errors.New("order has no items")

// discountRates maps a customer tier to a discount percentage.
var discountRates = map[string]float64{
	"gold":   0.20,
	"silver": 0.10,
}

// Process validates an order and returns its final total.
//
// Call chain: Process -> calculateTotal -> applyDiscount
func Process(o *Order) (float64, error) {
	if len(o.Items) == 0 {
		return 0, fmt.Errorf("process order %d: %w", o.ID, ErrEmptyOrder)
	}
	total := calculateTotal(o)
	return total, nil
}

func calculateTotal(o *Order) float64 {
	subtotal := 0.0
	for _, item := range o.Items {
		lineTotal := item.Price * float64(item.Quantity) // good breakpoint line
		subtotal += lineTotal
	}
	return applyDiscount(o.Customer, subtotal)
}

func applyDiscount(c *Customer, subtotal float64) float64 {
	// BUG (on purpose): no nil check on c. An order without a customer
	// panics here, which gives us a real stack trace to read.
	rate := discountRates[c.Tier]
	return subtotal * (1 - rate)
}
