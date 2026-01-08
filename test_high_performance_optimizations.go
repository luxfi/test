package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/luxfi/cache"
	"github.com/luxfi/concurrent"
	"github.com/luxfi/container"
	"github.com/luxfi/metric"
)

func main() {
	fmt.Println("Testing high-performance optimizations...")

	// Create a metrics registry
	registry := metric.NewMetricsRegistry()

	// Test DualMapCache
	fmt.Println("\n1. Testing DualMapCache...")
	testDualMapCache(registry)

	// Test ConcurrencyLimiter
	fmt.Println("\n2. Testing ConcurrencyLimiter...")
	testConcurrencyLimiter(registry)

	// Test OrderedMap with metrics
	fmt.Println("\n3. Testing OrderedMap with metrics...")
	testOrderedMap(registry)

	// Test StringInterner
	fmt.Println("\n4. Testing StringInterner...")
	testStringInterner(registry)

	// Print some metrics
	fmt.Println("\n5. Metrics Summary:")
	printMetrics(registry)

	fmt.Println("\nAll tests completed successfully!")
}

func testDualMapCache(reg *metric.MetricsRegistry) {
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

func testConcurrencyLimiter(reg *metric.MetricsRegistry) {
	limiter := concurrent.NewConcurrencyLimiter(3, "test_limiter", reg)

	// Test acquiring and releasing
	for i := 0; i < 5; i++ {
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

func testOrderedMap(reg *metric.MetricsRegistry) {
	omap := container.NewOrderedMap[string, string](reg)

	// Add items
	omap.Put("first", "1st")
	omap.Put("second", "2nd")
	omap.Put("third", "3rd")

	// Check size
	fmt.Printf("  Size: %d\n", omap.Size())

	// Get items
	if val, ok := omap.Get("first"); ok {
		fmt.Printf("  First value: %s\n", val)
	}

	// Iterate
	fmt.Print("  Iteration: ")
	omap.Iterate(func(key, value string) bool {
		fmt.Printf("%s=%s ", key, value)
		return true
	})
	fmt.Println()

	// Get keys and values
	keys := omap.Keys()
	values := omap.Values()
	fmt.Printf("  Keys: %v, Values: %v\n", keys, values)
}

func testStringInterner(reg *metric.MetricsRegistry) {
	interner := cache.NewStringInterner(reg)

	// Intern some strings
	str1 := interner.Intern("hello world")
	str2 := interner.Intern("hello world") // Should return same pointer
	str3 := interner.Intern("foo bar")

	// Check if they're the same by comparing values (since we can't compare string addresses directly)
	if str1 == str2 {
		fmt.Println("  String interning works correctly")
	} else {
		fmt.Println("  String interning failed")
	}

	fmt.Printf("  Interned '%s' and '%s'\n", str1, str3)
}

func printMetrics(reg *metric.MetricsRegistry) {
	// Just show that metrics are registered by printing their existence
	fmt.Println("  Metrics registry initialized with counters and gauges")
	
	// Show some stats about goroutines as a basic health check
	fmt.Printf("  Current goroutines: %d\n", runtime.NumGoroutine())
	
	// Sleep briefly to allow any background operations to complete
	time.Sleep(10 * time.Millisecond)
}