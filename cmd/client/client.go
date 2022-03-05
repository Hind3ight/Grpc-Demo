package main

import (
	"context"
	"fmt"
	pb "github.com/Hind3ight/Grpc-Demo/api/protocol"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)
	//runFirst(client)
	runSecond(client)

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

func runSecond(client pb.RouteGuideClient) {
	serverStream, err := client.ListFeatures(context.Background(), &pb.Rectangle{
		Lo: &pb.Point{
			Latitude:  23263493,
			Longitude: 113786833,
		},
		Hi: &pb.Point{
			Latitude:  29366260,
			Longitude: 113125491,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		feature, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(feature)
	}
}
