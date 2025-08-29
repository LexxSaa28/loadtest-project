package clients

import (
	"context"
	"fmt"
	"time"

	"load-test/stats"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	url   string
	stats stats.Collector
}

func NewWSClient(url string, stats stats.Collector) *WSClient {
	return &WSClient{
		url:   url,
		stats: stats,
	}
}

func (c *WSClient) MakeRequest(ctx context.Context, payloadSize int) (time.Duration, error) {
	// Establish connection
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, fmt.Sprintf("%s?size=%d", c.url, payloadSize), nil)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	
	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	
	start := time.Now()
	
	// Read message
	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	duration := time.Since(start)
	
	if err != nil {
		return duration, err
	}
	
	// Record metrics
	c.stats.Increment("requests.total")
	c.stats.Timing("request.time", duration)
	
	return duration, nil
}

func (c *WSClient) Close() error {
	return nil
}