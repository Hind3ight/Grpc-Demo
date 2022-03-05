package main

import (
	"context"
	"fmt"
	pb "github.com/Hind3ight/Grpc-Demo/api/protocol"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)
	runFirst(client)
}

func runFirst(client pb.RouteGuideClient) {
	features, err := client.GetFeatures(context.Background(), &pb.Point{
		Latitude:  33388590,
		Longitude: 120207474,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(features)
}
