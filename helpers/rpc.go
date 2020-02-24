package helpers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func RunGrpc(ip string, f func(context.Context, *grpc.ClientConn) (interface{}, error)) (interface{}, error) {
	log.Printf("Starting gRPC connection to %s", ip)
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("did not connect: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	ret, err := f(ctx, conn)

	cancel()
	_ = conn.Close()

	return ret, err
}
