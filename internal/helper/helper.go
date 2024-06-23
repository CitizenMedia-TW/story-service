package helper

import (
	"story-service/internal/database"
	"story-service/protobuffs/jwt-service"
)

type Helper struct {
	database   database.SQLDatabase
	JWTClient jwt.JWTServiceClient
}

func New(authClient jwt.JWTServiceClient) Helper {
	db := database.NewPostgresConn()

	// grpcClient, err := grpc.Dial("157.230.46.45:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	panic(err)
	// }

	return Helper{
		database:   db,
		JWTClient: authClient,
	}
}
