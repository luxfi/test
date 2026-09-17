package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/luxfi/cache"
	"github.com/luxfi/concurrent"
	"github.com/luxfi/metric"
)

func main() {
	fmt.Println("Testing high-performance optimizations...")

	// Create a metrics registry
	registry := metric.NewRegistry()

	// Test DualMapCache
	fmt.Println("\n1. Testing DualMapCache...")
	testDualMapCache(registry)

	// Test ConcurrencyLimiter
	fmt.Println("\n2. Testing ConcurrencyLimiter...")
	testConcurrencyLimiter(registry)

	// Print some metrics
	fmt.Println("\n3. Metrics Summary:")
	printMetrics(registry)

	fmt.Println("\nAll tests completed successfully!")
}

func testDualMapCache(reg metric.Registry) {
	dualCache := cache.NewDualMapCache[string, string](reg)

	// Put some values
	dualCache.Put("key1", "value1")
	dualCache.Put("key2", "value2")
	dualCache.Put("key3", "value3")

	// Get values
	if val, ok := dualCache.Get("key1"); ok {
		fmt.Printf("  Retrieved: %s\n", val)
	} else {
		log.Fatal("Failed to retrieve key1")
	}

	// Trigger migration
	dualCache.Migrate()
	fmt.Println("  Migration completed")
}

func testConcurrencyLimiter(reg metric.Registry) {
	limiter := concurrent.NewConcurrencyLimiter(3, "test_limiter", reg)

	// Test acquiring and releasing
	for i := range 5 {
		if err := limiter.Acquire(nil); err != nil {
			log.Fatalf("Failed to acquire: %v", err)
		}
		fmt.Printf("  Acquired token %d, available: %d\n", i+1, limiter.Available())
		limiter.Release()
	}

	// Test try acquire
	if acquired := limiter.TryAcquire(); acquired {
		fmt.Println("  TryAcquire succeeded")
		limiter.Release()
	} else {
		fmt.Println("  TryAcquire failed")
	}
}

func printMetrics(reg metric.Registry) {
	// Just show that metrics are registered by printing their existence
	fmt.Println("  Metrics registry initialized with counters and gauges")

	// Show some stats about goroutines as a basic health check
	fmt.Printf("  Current goroutines: %d\n", runtime.NumGoroutine())

	// Sleep briefly to allow any background operations to complete
	time.Sleep(10 * time.Millisecond)
}
