// Package utils provides utility types and functions for the limit order book.
package utils

// Import necessary packages.
import (
	"fmt"                             // For formatted I/O operations.
	"github.com/google/uuid"          // For generating unique identifiers for orders.
	"github.com/shopspring/decimal"   // For high-precision decimal arithmetic operations.
	"time"                            // For working with time values.
)

// LimitOrder represents an individual limit order in the order book.
type LimitOrder struct {
	ID        uuid.UUID       // Unique identifier for the order.
	Ticker    string          // Stock or asset ticker symbol associated with the order.
	Price     decimal.Decimal // Price at which the order is placed, using high-precision decimals.
	Quantity  int             // Number of shares (or units) specified in the order.
	IsBid     bool            // Boolean flag indicating if the order is a bid (buy) order; false implies an ask (sell) order.
	Timestamp time.Time       // The time at which the order was created.
}

// String implements the Stringer interface for LimitOrder.
// It returns a formatted string representation of the LimitOrder.
func (lo LimitOrder) String() string {
	// Format the LimitOrder details into a string.
	// Format: [Ticker - Price - Quantity - Is Buy: IsBid - Timestamp in milliseconds - Order ID]
	return fmt.Sprintf("[%s - %s - %d - Is Buy: %t - %d - %s]\n",
		lo.Ticker,           // Ticker symbol.
		lo.Price.String(),   // Price formatted as a string.
		lo.Quantity,         // Order quantity.
		lo.IsBid,            // Boolean indicating if the order is a bid.
		lo.Timestamp.UnixMilli(), // Timestamp converted to milliseconds.
		lo.ID.String())      // Unique order ID as a string.
}

// OrderPriorityQueue is a slice of pointers to LimitOrder that implements heap.Interface.
// It is used as a priority queue for managing orders.
type OrderPriorityQueue []*LimitOrder

// Len returns the number of orders in the priority queue.
// This method is required by the heap.Interface.
func (pq OrderPriorityQueue) Len() int {
	return len(pq) // Return the length of the slice.
}

// Less compares two orders in the priority queue based on their prices.
// For bid orders, a higher price has higher priority (max-heap behavior).
// For ask orders, a lower price has higher priority (min-heap behavior).
func (pq OrderPriorityQueue) Less(i, j int) bool {
	// Check if the order at index i is a bid order.
	if pq[i].IsBid {
		// For bid orders, return true if the price at index i is greater than the price at index j.
		return pq[i].Price.GreaterThan(pq[j].Price)
	}
	// For ask orders, return true if the price at index i is less than the price at index j.
	return pq[i].Price.LessThan(pq[j].Price)
}

// Swap exchanges the orders at indices i and j in the priority queue.
// This method is required by the heap.Interface.
func (pq OrderPriorityQueue) Swap(i, j int) {
	// Swap the two orders.
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds a new order to the priority queue.
// x is expected to be a pointer to a LimitOrder.
func (pq *OrderPriorityQueue) Push(x any) {
	// Convert x to a pointer to LimitOrder.
	limitOrder := x.(*LimitOrder)
	// Append the new order to the end of the slice.
	*pq = append(*pq, limitOrder)
}

// Pop removes and returns the last order from the priority queue.
// This method is required by the heap.Interface.
func (pq *OrderPriorityQueue) Pop() any {
	// Create a temporary variable to hold the current slice.
	old := *pq
	// Get the current length of the slice.
	n := len(old)
	// Get the last element (order) in the slice.
	limitOrder := old[n-1]
	// Remove the last element from the slice.
	*pq = old[0 : n-1]
	// Return the removed order.
	return limitOrder
}

// Peek returns the top order of the priority queue without removing it.
// This is useful for inspecting the highest priority order.
func (pq *OrderPriorityQueue) Peek() any {
	// Retrieve the current slice.
	old := *pq
	// Return the first element in the slice, which is the highest priority order.
	return old[0]
}
