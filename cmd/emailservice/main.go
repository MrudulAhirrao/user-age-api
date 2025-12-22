package main

import (
	"net"
	"log"
	"google.golang.org/grpc"
	"user-age-api/internal/email"
	pb "user-age-api/internal/email/proto"
)

func main(){
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Email Sevrice Running on Port %s",port)
	grpc := grpc.NewServer()
	emailServer:= &email.Server{}
	pb.RegisterEmailServiceServer(grpc, emailServer)
	if err := grpc.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}