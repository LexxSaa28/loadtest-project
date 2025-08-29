package main

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	pb "grpc-server/proto"

	"google.golang.org/grpc"
)

const (
	ALPHANUMERIC_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DEFAULT_PAYLOAD_SIZE_KB = 1
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func generateRandomPayload(sizeKB int) string {
	totalChars := sizeKB * 1024
	var builder strings.Builder
	builder.Grow(totalChars)

	for i := 0; i < totalChars; i++ {
		randomIndex := rand.Intn(len(ALPHANUMERIC_CHARS))
		builder.WriteByte(ALPHANUMERIC_CHARS[randomIndex])
	}

	return builder.String()
}

func (s *server) SayHello(req *pb.HelloRequest, stream pb.HelloService_SayHelloServer) error {
	sizeKB := DEFAULT_PAYLOAD_SIZE_KB
	if req.SizeKb > 0 {
		sizeKB = int(req.SizeKb)
	}

	payload := generateRandomPayload(sizeKB)
	timestamp := time.Now().Format(time.RFC3339Nano)

	response := &pb.HelloResponse{
		Message:       "Hello from gRPC backend!",
		Timestamp:     timestamp,
		PayloadSizeKb: int32(sizeKB),
		Payload:       payload,
	}

	return stream.Send(response)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	
	log.Printf("gRPC server running on :9001")
	log.Printf("Default payload size: %d KB", DEFAULT_PAYLOAD_SIZE_KB)
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}