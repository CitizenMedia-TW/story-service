package app

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"story-service/internal/restapp"
	"story-service/protobuffs/jwt-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func StartServer() {
	grpcClient, err := grpc.Dial("157.230.46.45:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	jwtClient := jwt.NewJWTServiceClient(grpcClient)

	restServer := restapp.New(jwtClient)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = http.Serve(lis, restServer.Routes())
	if err != nil {
		return
	}
}
