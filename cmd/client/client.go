package main

import (
	"context"
	"fmt"
	pb "github.com/Hind3ight/Grpc-Demo/api/protocol"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewRouteGuideClient(conn)
	//runFirst(client)
	//runSecond(client)
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

func runThird(client pb.RouteGuideClient) {
	points := []*pb.Point{
		{
			Latitude: 22287502, Longitude: 114149268,
		},
		{
			Latitude: 22480666, Longitude: 113418481,
		},
		{
			Latitude: 40068078, Longitude: 124330737,
		},
	}

	clientStream, err := client.RecordRoute(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	for _, point := range points {
		if err := clientStream.Send(point); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}

	summary, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(summary)
}
