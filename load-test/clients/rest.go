package clients

import (
	"context"
	"fmt"
	"net/http"
	"time"
	
	"load-test/stats"
)

type RESTClient struct {
	baseURL string
	client  *http.Client
	stats   stats.Collector
}

func NewRESTClient(baseURL string, stats stats.Collector) *RESTClient {
	return &RESTClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
		stats:   stats,
	}
}

func (c *RESTClient) MakeRequest(ctx context.Context, payloadSize int) (time.Duration, error) {
	url := fmt.Sprintf("%s?size=%d", c.baseURL, payloadSize)
	
	start := time.Now()
	resp, err := c.client.Get(url)
	duration := time.Since(start)
	
	if err != nil {
		return duration, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return duration, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	// Record metrics
	c.stats.Increment("requests.total")
	c.stats.Timing("request.time", duration)
	
	return duration, nil
}

func (c *RESTClient) Close() error {
	return nil
}