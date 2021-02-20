package main

import (
	"google.golang.org/grpc/reflection"
	"github.com/reseich/FullCycle-GRPC/pb"
	"github.com/reseich/FullCycle-GRPC/services"
	"net" 
	"log"
	"google.golang.org/grpc"
)

func main()  {
	lis, err := net.Listen("tcp", "localhost:50051")
	if lis != nil {
		log.Print("connection opened in  localhost:50051 ")
	}

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}	

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, services.NewUserService())
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Could not serve: %v", err)
	}

}
