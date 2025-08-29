package clients

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type LoadTestResult struct {
	TotalRequests   int
	Successful      int
	Failed          int
	TotalDuration   time.Duration
	MinDuration     time.Duration
	MaxDuration     time.Duration
	AvgDuration     time.Duration
	RequestsPerSecond float64
}

func (r *LoadTestResult) Print() {
	fmt.Printf("Load Test Results:\n")
	fmt.Printf("  Total Requests: %d\n", r.TotalRequests)
	fmt.Printf("  Successful: %d\n", r.Successful)
	fmt.Printf("  Failed: %d\n", r.Failed)
	fmt.Printf("  Total Duration: %v\n", r.TotalDuration)
	fmt.Printf("  Min Duration: %v\n", r.MinDuration)
	fmt.Printf("  Max Duration: %v\n", r.MaxDuration)
	fmt.Printf("  Avg Duration: %v\n", r.AvgDuration)
	fmt.Printf("  Requests Per Second: %.2f\n", r.RequestsPerSecond)
}

type LoadTestRunner struct {
	client      Client
	requests    int
	concurrency int
	payloadSize int
	duration    time.Duration
}

func NewLoadTestRunner(client Client, requests, concurrency, payloadSize int, duration time.Duration) *LoadTestRunner {
	return &LoadTestRunner{
		client:      client,
		requests:    requests,
		concurrency: concurrency,
		payloadSize: payloadSize,
		duration:    duration,
	}
}

func (r *LoadTestRunner) Run() *LoadTestResult {
	var wg sync.WaitGroup
	results := make(chan time.Duration, r.requests)
	errors := make(chan error, r.requests)
	
	ctx, cancel := context.WithTimeout(context.Background(), r.duration)
	defer cancel()
	
	startTime := time.Now()
	
	// Start workers
	for i := 0; i < r.concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for ctx.Err() == nil {
				duration, err := r.client.MakeRequest(ctx, r.payloadSize)
				if err != nil {
					errors <- err
				} else {
					results <- duration
				}
			}
		}(i)
	}
	
	// Wait for all workers to finish or timeout
	wg.Wait()
	close(results)
	close(errors)
	
	totalDuration := time.Since(startTime)
	
	// Process results
	result := &LoadTestResult{
		TotalRequests: r.requests,
		TotalDuration: totalDuration,
		MinDuration:   time.Hour, // Initialize with a large value
	}
	
	// Count successful requests and calculate durations
	var totalRequestTime time.Duration
	count := 0
	
	for duration := range results {
		count++
		totalRequestTime += duration
		
		if duration < result.MinDuration {
			result.MinDuration = duration
		}
		if duration > result.MaxDuration {
			result.MaxDuration = duration
		}
	}
	
	// Count errors
	errorCount := 0
	for range errors {
		errorCount++
	}
	
	result.Successful = count
	result.Failed = errorCount
	result.AvgDuration = totalRequestTime / time.Duration(count)
	result.RequestsPerSecond = float64(count) / totalDuration.Seconds()
	
	return result
}