// File: internal/grpc/aaa_client.go
package grpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb_aaa "github.com/kisanlink/protobuf/pb-aaa"
)

// AAAGrpcClient wraps the gRPC client for AAA service
type AAAGrpcClient struct {
	conn   *grpc.ClientConn
	client pb_aaa.UserServiceClient
}

// NewAAAGrpcClient creates a new gRPC client for AAA service
func NewAAAGrpcClient(endpoint string) (*AAAGrpcClient, error) {
	log.Printf("🔌 Connecting to AAA gRPC service at: %s", endpoint)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AAA gRPC service: %w", err)
	}

	client := pb_aaa.NewUserServiceClient(conn)
	log.Printf("✅ Successfully connected to AAA gRPC service")

	return &AAAGrpcClient{conn: conn, client: client}, nil
}

// Close closes the gRPC connection
func (c *AAAGrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUserByMobileNumber retrieves a user by mobile number from AAA service.
// It expects the proto field MobileNumber to be uint64.
func (c *AAAGrpcClient) GetUserByMobileNumber(ctx context.Context, mobile string) (*pb_aaa.GetUserByMobileNumberResponse, error) {
	log.Printf("📱 AAA gRPC GetUserByMobileNumber - Mobile: %s", mobile)

	// Convert string to uint64 as per proto definition
	mobileUint, err := strconv.ParseUint(mobile, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid mobile number format: %w", err)
	}

	req := &pb_aaa.GetUserByMobileNumberRequest{
		MobileNumber: mobileUint,
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := c.client.GetUserByMobileNumber(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AAA service GetUserByMobileNumber failed: %w", err)
	}

	log.Printf("✅ AAA gRPC GetUserByMobileNumber successful: %s", resp.GetMessage())
	return resp, nil
}
