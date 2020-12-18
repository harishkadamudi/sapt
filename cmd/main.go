package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ishankhare07/sapt/pkg/logger"
	"google.golang.org/grpc"

	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	testv3 "github.com/envoyproxy/go-control-plane/pkg/test/v3"
	"github.com/ishankhare07/sapt/pkg/xds"
)

var port uint

func init() {
	port = 18000
}

func main() {
	l := logger.Logger{Debug: true}
	cache := cachev3.NewSnapshotCache(false, cachev3.IDHash{}, l)
	cb := &testv3.Callbacks{Debug: true}
	server := serverv3.NewServer(context.Background(), cache, cb)

	grpcServer := grpc.NewServer()
	xds.RegisterXDSHandlers(grpcServer, server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	log.Printf("server listeneing on port: %d", port)
	if err = grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
