package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"example.com/go-debug-demo/internal/order"
)

func main() {
	crash := flag.Bool("crash", false, "process an order with a nil customer to trigger a panic")
	trace := flag.Bool("trace", true, "print the current stack trace from inside the program")
	flag.Parse()

	orders := []*order.Order{
		{
			ID:       1,
			Customer: &order.Customer{Name: "Alice", Tier: "gold"},
			Items: []order.Item{
				{SKU: "KB-01", Price: 49.90, Quantity: 1},
				{SKU: "MS-02", Price: 19.50, Quantity: 2},
			},
		},
		{
			ID:       2,
			Customer: &order.Customer{Name: "Bob", Tier: "silver"},
			Items:    nil, // returns ErrEmptyOrder
		},
	}

	if *crash {
		orders = append(orders, &order.Order{
			ID:       3,
			Customer: nil, // triggers the nil-pointer panic in applyDiscount
			Items:    []order.Item{{SKU: "HD-03", Price: 99.00, Quantity: 1}},
		})
	}

	if *trace {
		printStack()
	}

	for _, o := range orders {
		total, err := order.Process(o)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			continue
		}
		fmt.Printf("order %d total: %.2f\n", o.ID, total)
	}
}

// printStack shows two ways to capture a stack trace programmatically.
func printStack() {
	fmt.Println("=== debug.PrintStack() ===")
	debug.PrintStack()

	fmt.Println("\n=== runtime.Callers() frame by frame ===")
	pcs := make([]uintptr, 16)
	n := runtime.Callers(1, pcs) // 1 = skip runtime.Callers itself
	frames := runtime.CallersFrames(pcs[:n])
	for {
		f, more := frames.Next()
		fmt.Printf("%s\n\t%s:%d\n", f.Function, f.File, f.Line)
		if !more {
			break
		}
	}
	fmt.Println()
}
