// Package main contains the entry point for the executable.
package main

// Import necessary packages.
import (
	"fmt"       // Package for formatted I/O operations.
	"math/rand" // Package to generate pseudo-random numbers.
	"strconv"   // Package for conversions to and from string representations.
	"sync"      // Package that provides basic synchronization primitives such as mutual exclusion locks.
	"time"      // Package for measuring and displaying time.

	"github.com/manuelinfosec/limit-orderbook-go/entities" // Custom package containing OrderBook implementation.
)

// Constants for the simulation.
const (
	N            = 1_000_000 // Total number of orders to add for each order type per ticker.
	BuyMaxPrice  = 300       // Maximum price for buy (bid) orders.
	BuyMinPrice  = 150       // Minimum price for buy (bid) orders.
	SellMaxPrice = 450       // Maximum price for sell (ask) orders.
	SellMinPrice = 250       // Minimum price for sell (ask) orders.
)

// main is the entry point of the application.
func main() {
	// Create a map of ticker symbols to their respective OrderBook instances.
	orderbooks := map[string]*entities.OrderBook{
		"AAPL": entities.NewOrderBook("AAPL"), // Initialize an order book for ticker "AAPL".
		"HOOD": entities.NewOrderBook("HOOD"), // Initialize an order book for ticker "HOOD".
		"SPY":  entities.NewOrderBook("SPY"),  // Initialize an order book for ticker "SPY".
		"SHOP": entities.NewOrderBook("SHOP"), // Initialize an order book for ticker "SHOP".
		"QQQ":  entities.NewOrderBook("QQQ"),  // Initialize an order book for ticker "QQQ".
	}

	// Create a WaitGroup to synchronize concurrent goroutines.
	wg := &sync.WaitGroup{}

	// Print a message indicating the start of limit order population.
	fmt.Println("Populating limit orders")
	// Record the start time to later measure how long the population process takes.
	orderPopulateStart := time.Now()

	// Loop over each ticker and corresponding order book to add bid orders concurrently.
	for ticker, ob := range orderbooks {
		wg.Add(1)                                  // Increment the WaitGroup counter for each new goroutine.
		go addTestLimitOrder(wg, ob, ticker, true) // Launch a goroutine to add buy (bid) orders.
	}

	// Loop over each ticker and corresponding order book to add ask orders concurrently.
	for ticker, ob := range orderbooks {
		wg.Add(1)                                   // Increment the WaitGroup counter for each new goroutine.
		go addTestLimitOrder(wg, ob, ticker, false) // Launch a goroutine to add sell (ask) orders.
	}

	// Wait for all order population goroutines to complete.
	wg.Wait()
	// Calculate and print the total time taken to populate limit orders (in seconds).
	fmt.Printf("Time to populate limit orders: %d seconds", (time.Now().UnixMilli()-orderPopulateStart.UnixMilli())/int64(1000))

	// Loop over each ticker to print the total count of bid and ask orders in the respective order books.
	for ticker, ob := range orderbooks {
		fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	}

	// Record the start time for the order matching process.
	orderMatchingStart := time.Now()
	// Print a message indicating the beginning of the order matching process.
	fmt.Println("Order matching begins...")
	// Launch goroutines for order matching concurrently for each ticker.
	for ticker, ob := range orderbooks {
		wg.Add(1)                          // Increment the WaitGroup counter for each matching process goroutine.
		go executeMatching(ob, ticker, wg) // Launch a goroutine to execute matching on the order book.
	}

	// Wait for all matching execution goroutines to complete.
	wg.Wait()
	// Calculate and print the total time taken to finish order matching (in seconds).
	fmt.Printf("Time to finish order matching: %d seconds", (time.Now().UnixMilli()-orderMatchingStart.UnixMilli())/int64(1000))
}

// executeMatching triggers the matching algorithm on a given order book.
// It then prints the remaining number of bid and ask orders for the ticker.
// Finally, it decrements the WaitGroup counter to signal completion.
func executeMatching(ob *entities.OrderBook, ticker string, wg *sync.WaitGroup) {
	ob.Match() // Invoke the Match method to process order matching within the order book.
	// Print the ticker symbol along with the number of remaining bid and ask orders.
	fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	wg.Done() // Signal that this goroutine's work is complete.
}

// addTestLimitOrder adds N limit orders to the provided order book.
// It generates random prices and quantities, and designates the orders as bid or ask based on isBid.
func addTestLimitOrder(wg *sync.WaitGroup, ob *entities.OrderBook, ticker string, isBid bool) {
	// Loop N times to add the predetermined number of orders.
	for i := 0; i < N; i++ {
		var randomPrice float64 // Declare a variable to hold the generated random price.
		// Create a new random number generator with a seed based on the current time in nanoseconds.
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		if isBid {
			// For buy orders, generate a random price between BuyMinPrice and BuyMaxPrice.
			randomPrice = BuyMinPrice + r.Float64()*(BuyMaxPrice-BuyMinPrice)
		} else {
			// For sell orders, generate a random price between SellMinPrice and SellMaxPrice.
			randomPrice = SellMinPrice + r.Float64()*(SellMaxPrice-SellMinPrice)
		}

		// Format the random price as a string with 2 decimal places.
		priceString := strconv.FormatFloat(randomPrice, 'f', 2, 64)
		// Generate a random quantity for the order, with a maximum possible value of 50,000.
		quantity := r.Intn(50000)

		// Add the generated limit order to the order book with the provided ticker,
		// price string, quantity, and order type (bid or ask).
		ob.AddLimitOrder(
			ticker,      // Ticker symbol associated with the order.
			priceString, // Price of the order formatted as a string.
			quantity,    // Quantity of the order.
			isBid)       // Boolean flag indicating if the order is a bid (true) or an ask (false).

		// The following print statement is commented out.
		// It can be enabled for debugging to track each created limit order.
		// fmt.Printf("Created limit order - ticker: %s price: %s quantity: %d isBid: %t \n", ticker, priceString, quantity, isBid)
	}

	// Signal that this goroutine has finished adding orders by decrementing the WaitGroup counter.
	wg.Done()
}
