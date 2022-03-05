package main

import (
	"context"
	"encoding/json"
	pb "github.com/Hind3ight/Grpc-Demo/api/protocol" // 引入route.proto的package
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net"
)

type routeGuideServer struct {
	features                         []*pb.Feature
	pb.UnimplementedRouteGuideServer // 需内嵌该服务，因为protoc解析.proto文件并创建代码桩时，会用该服务完成向前兼容的实现
}

func (s *routeGuideServer) GetFeatures(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.features {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalln("cannot create a listener at the address")
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	log.Fatalln(grpcServer.Serve(lis))
}

func newServer() *routeGuideServer {
	return &routeGuideServer{
		features: getSource(),
	}
}

func getSource() []*pb.Feature {
	data, err := ioutil.ReadFile("./sources/features.txt")
	if err != nil {
		log.Fatal(err)
	}
	var features []*pb.Feature
	err = json.Unmarshal(data, &features)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return features
}
