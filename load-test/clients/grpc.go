package clients

import (
	"context"
	"time"

	pb "grpc-server/proto"
	"load-test/stats"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn  *grpc.ClientConn
	client pb.HelloServiceClient
	stats stats.Collector
}

func NewGRPCClient(addr string, stats stats.Collector) *GRPCClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	
	return &GRPCClient{
		conn:   conn,
		client: pb.NewHelloServiceClient(conn),
		stats:  stats,
	}
}

func (c *GRPCClient) MakeRequest(ctx context.Context, payloadSize int) (time.Duration, error) {
	req := &pb.HelloRequest{SizeKb: int32(payloadSize)}
	
	start := time.Now()
	stream, err := c.client.SayHello(ctx, req)
	if err != nil {
		return time.Since(start), err
	}
	
	// Receive response
	_, err = stream.Recv()
	duration := time.Since(start)
	
	if err != nil {
		return duration, err
	}
	
	// Record metrics
	c.stats.Increment("requests.total")
	c.stats.Timing("request.time", duration)
	
	return duration, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}