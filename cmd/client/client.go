package main

import (
	"time"
	"context"
	"github.com/reseich/FullCycle-GRPC/pb"
	"log"
	"io"
	"google.golang.org/grpc"
)

func main()  {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err !=nil{
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	log.Print("------------------------------Unary WAY---------------------------")
	AddUser(client)
	log.Print("------------------------------Server Stream WAY---------------------------")
	AddUserVerbose(client)
	log.Print("------------------------------Client Stream WAY---------------------------")
	AddUsers(client)
	log.Print("------------------------------Bidirectional Stream WAY---------------------------")
	AddUsersStreamBoth(client)
}

func AddUser(client pb.UserServiceClient){
	req := &pb.User{
		Id: "0",
		Name: "Rafael",
		Email: "reseich@gmail.com",
	}
	res, err := client.AddUser(context.Background(), req)
	
	if err !=nil{
		log.Fatalf("Could not make gRPC Server: %v", err)
	}

	log.Print(res)
}

func AddUserVerbose(client pb.UserServiceClient){
	req := &pb.User{
		Id: "0",
		Name: "Rafael",
		Email: "reseich@gmail.com",
	}
	responseStream, err := client.AddUserVerbose(context.Background(), req)
	
	if err != nil{
		log.Fatalf("Error: %v", err)
	}

	for{
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break;
		}

		if err !=nil{
			log.Fatalf("Could not receive the message: %v", err)
		}
		log.Print(stream.Status)
	}
	
}

func AddUsers(client pb.UserServiceClient){

	reqs := []*pb.User{
		&pb.User{
			Id: "r0",
			Name: "Rafael",
			Email: "reseich@gmail.com",
		},
		&pb.User{
			Id: "r1",
			Name: "Rafael1",
			Email: "reseich1@gmail.com",
		},
		&pb.User{
			Id: "r2",
			Name: "Rafael2",
			Email: "reseich2@gmail.com",
		},
		&pb.User{
			Id: "r3",
			Name: "Rafael3",
			Email: "reseich3@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for  _, req := range reqs {
		stream.Send(req)
		log.Print("sending ", req.Name)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving request: %v", err)
	}

	log.Print(res)

}

func AddUsersStreamBoth(client pb.UserServiceClient){
	
	reqs := []*pb.User{
		&pb.User{
			Id: "r0",
			Name: "Rafael",
			Email: "reseich@gmail.com",
		},
		&pb.User{
			Id: "r1",
			Name: "Rafael1",
			Email: "reseich1@gmail.com",
		},
		&pb.User{
			Id: "r2",
			Name: "Rafael2",
			Email: "reseich2@gmail.com",
		},
		&pb.User{
			Id: "r3",
			Name: "Rafael3",
			Email: "reseich3@gmail.com",
		},
	}
	

	stream, err := client.AddUsersStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	
	wait := make(chan int)

	//routine to sending all users
	go func(){
		for  _, req := range reqs {
			log.Print("sending user ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 3)
		}
		stream.CloseSend()
	}()

	//routine to receive all users
	go func(){
		for{
			res, err := stream.Recv()
			if err == io.EOF {
				break;
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}
			log.Print("receiving user ", res.GetUser().GetName()," with status ",res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}