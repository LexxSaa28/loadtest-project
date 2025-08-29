package main

import (
	"context"
	"log"
	"os"
	"time"

	"load-test/clients"
	"load-test/stats"

	"github.com/spf13/cobra"
)

func main() {
	var protocol string
	var requests int
	var concurrency int
	var payloadSize int
	var duration time.Duration

	var rootCmd = &cobra.Command{
		Use:   "loadtest",
		Short: "Load testing tool for multiple protocols",
		Run: func(cmd *cobra.Command, args []string) {
			// Initialize stats collector
			statsCollector := stats.NewStatsDCollector("localhost:8125")
			
			var client clients.Client
			
			switch protocol {
			case "rest":
				client = clients.NewRESTClient("http://localhost:9000/hello", statsCollector)
			case "grpc":
				client = clients.NewGRPCClient("localhost:9001", statsCollector)
			case "ws":
				client = clients.NewWSClient("ws://localhost:9002/ws", statsCollector)
			default:
				log.Fatalf("Unknown protocol: %s", protocol)
			}
			
			// Run load test
			runner := clients.NewLoadTestRunner(client, requests, concurrency, payloadSize, duration)
			results := runner.Run()
			
			// Print results
			results.Print()
		},
	}

	rootCmd.Flags().StringVarP(&protocol, "protocol", "p", "rest", "Protocol to test (rest|grpc|ws)")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1000, "Total number of requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent workers")
	rootCmd.Flags().IntVarP(&payloadSize, "size", "s", 1, "Payload size in KB")
	rootCmd.Flags().DurationVarP(&duration, "duration", "d", 10*time.Second, "Test duration")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}