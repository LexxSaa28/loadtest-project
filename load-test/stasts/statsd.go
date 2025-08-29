package stats

import (
	"fmt"
	"log"
	"net"
	"time"
)

type StatsDCollector struct {
	addr   string
	prefix string
	conn   net.Conn
}

func NewStatsDCollector(addr string) *StatsDCollector {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Printf("Failed to connect to StatsD: %v", err)
		return &StatsDCollector{} // Return dummy collector
	}
	
	return &StatsDCollector{
		addr:   addr,
		prefix: "loadtest",
		conn:   conn,
	}
}

func (c *StatsDCollector) Increment(metric string) {
	if c.conn == nil {
		return
	}
	
	msg := fmt.Sprintf("%s.%s:1|c", c.prefix, metric)
	c.conn.Write([]byte(msg))
}

func (c *StatsDCollector) Timing(metric string, duration time.Duration) {
	if c.conn == nil {
		return
	}
	
	msg := fmt.Sprintf("%s.%s:%d|ms", c.prefix, metric, duration.Milliseconds())
	c.conn.Write([]byte(msg))
}

func (c *StatsDCollector) Gauge(metric string, value float64) {
	if c.conn == nil {
		return
	}
	
	msg := fmt.Sprintf("%s.%s:%.2f|g", c.prefix, metric, value)
	c.conn.Write([]byte(msg))
}

func (c *StatsDCollector) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}