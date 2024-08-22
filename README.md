# limit-orderbook-go

A Golang implementation of a Level 2 order book and matching engine.  

The program initializes five independent order book instances and populates each with one million limit orders per order type per ticker using goroutines—resulting in a total of 10 million orders. Once populated, it concurrently executes order matching across all tickers.  

Execution time is approximately 45 seconds, with 25 seconds allocated to order insertion and 20 seconds to matching execution. Performance benchmarks on an M1 Pro CPU indicate speeds of around 400K order insertions per second and 500K order fulfillments per second.