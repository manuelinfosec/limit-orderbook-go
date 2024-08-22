// Package entities contains the core structures and methods for the order book.
package entities

// Import standard and external packages.
import (
	"container/heap" // Provides heap operations for implementing priority queues.
	"fmt"            // Package for formatted I/O.
	"log"            // Package for logging errors and information.
	"time"           // Package for time-related functions.

	"github.com/google/uuid"                                  // Package to generate unique identifiers (UUIDs) for orders.
	utils "github.com/manuelinfosec/limit-orderbook-go/utils" // Importing the utils package with alias 'utils' which contains definitions for LimitOrder and OrderPriorityQueue.
	"github.com/shopspring/decimal"                           // Package for high-precision decimal arithmetic.
)

// OrderBook struct represents a market order book for a specific ticker.
type OrderBook struct {
	ticker string                   // The ticker symbol that this order book is associated with.
	Bids   utils.OrderPriorityQueue // A priority queue (heap) that stores bid (buy) orders.
	Asks   utils.OrderPriorityQueue // A priority queue (heap) that stores ask (sell) orders.
}

// NewOrderBook returns a new instance of OrderBook for a given ticker.
func NewOrderBook(ticker string) *OrderBook {
	// Create a new OrderBook instance with empty bid and ask queues.
	ob := &OrderBook{
		ticker: ticker,                            // Set the ticker symbol for the order book.
		Bids:   make(utils.OrderPriorityQueue, 0), // Initialize the bid queue with zero length.
		Asks:   make(utils.OrderPriorityQueue, 0), // Initialize the ask queue with zero length.
	}

	// Initialize the bid priority queue as a heap.
	heap.Init(&ob.Bids)
	// Initialize the ask priority queue as a heap.
	heap.Init(&ob.Asks)
	// Return the pointer to the newly created OrderBook.
	return ob
}

// intMin returns the smaller of two integer values.
func intMin(a int, b int) int {
	if a < b { // Check if a is less than b.
		return a // Return a if it is the smaller value.
	}

	return b // Otherwise, return b.
}

// AddLimitOrder adds a new limit order to the order book.
// It takes the ticker, price as a string, quantity, and a boolean indicating if it's a bid order.
func (ob *OrderBook) AddLimitOrder(ticker string, priceString string, quantity int, IsBid bool) bool {
	// Validate that the provided ticker matches the order book's ticker.
	if ticker != ob.ticker {
		// Log a fatal error if the tickers do not match and terminate execution.
		log.Fatalf("Wrong ticker to place a new order. Unable to create a new order: %s, %s, %d, IsBid: %t", ticker, priceString, quantity, IsBid)
		return false // This return is unreachable due to log.Fatalf, but it's included to satisfy the function signature.
	}

	// Convert the price string to a decimal type for precise arithmetic.
	priceDecimal, err := decimal.NewFromString(priceString)
	// If an error occurs during conversion, log the error and exit.
	if err != nil {
		log.Fatalf("Invalid price found. Unable to create a new order: %s, %s, %d, IsBid: %t", ticker, priceString, quantity, IsBid)
		return false // This return is unreachable because log.Fatalf exits, but it's provided as a safeguard.
	}

	// Create a new LimitOrder struct with the provided values.
	newOrder := &utils.LimitOrder{
		IsBid:     IsBid,        // Set whether the order is a bid (buy) order.
		ID:        uuid.New(),   // Generate a new unique identifier for the order.
		Price:     priceDecimal, // Set the price of the order as a decimal.
		Quantity:  quantity,     // Set the quantity for the order.
		Ticker:    ticker,       // Set the ticker symbol for the order.
		Timestamp: time.Now(),   // Record the current time as the order's timestamp.
	}

	// Depending on whether the order is a bid or an ask, push it into the corresponding heap.
	if newOrder.IsBid {
		heap.Push(&ob.Bids, newOrder) // For bid orders, push onto the Bids heap.
	} else {
		heap.Push(&ob.Asks, newOrder) // For ask orders, push onto the Asks heap.
	}

	// Return true indicating the order was added successfully.
	return true
}

// Match processes the order matching within the order book.
func (ob *OrderBook) Match() {
	// Continue matching as long as there are orders in both the bid and ask queues.
	for ob.Bids.Len() > 0 && ob.Asks.Len() > 0 {
		// Peek at the highest priority bid order without removing it from the heap.
		buy := ob.Bids.Peek().(*utils.LimitOrder)
		// Peek at the highest priority ask order without removing it from the heap.
		sell := ob.Asks.Peek().(*utils.LimitOrder)

		// Print details of the bid and ask orders that are being considered for matching.
		fmt.Printf("bid/ask to match: %d shares at %s VS %d shares at %s\n", buy.Quantity, buy.Price.String(), sell.Quantity, sell.Price.String())

		// Check if the bid price is lower than the ask price.
		// If true, no further matching can occur because the highest bid is less than the lowest ask.
		if buy.Price.LessThan(sell.Price) {
			break // Exit the loop since no orders can be matched under the current conditions.
		}

		// Remove the highest priority bid order from the heap.
		buy = heap.Pop(&ob.Bids).(*utils.LimitOrder)
		// Remove the highest priority ask order from the heap.
		sell = heap.Pop(&ob.Asks).(*utils.LimitOrder)

		// Determine the number of shares to be traded by finding the minimum of the two order quantities.
		quantityFilled := intMin(buy.Quantity, sell.Quantity)
		// Print the trade execution details: ticker, number of shares matched, and the trade price.
		fmt.Printf("Ticker %s - Matched %d shares at %s\n", ob.ticker, quantityFilled, sell.Price.String())

		// Deduct the matched quantity from the bid order.
		buy.Quantity -= quantityFilled
		// Deduct the matched quantity from the ask order.
		sell.Quantity -= quantityFilled

		// If the bid order still has remaining shares (i.e., partially filled), push it back onto the heap.
		if buy.Quantity > 0 {
			heap.Push(&ob.Bids, buy)
		}

		// If the ask order still has remaining shares (i.e., partially filled), push it back onto the heap.
		if sell.Quantity > 0 {
			heap.Push(&ob.Asks, sell)
		}
	}
}
