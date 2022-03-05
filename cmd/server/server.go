package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Hind3ight/Grpc-Demo/api/protocol" // 引入route.proto的package
	"github.com/Hind3ight/Grpc-Demo/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"
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

func (s *routeGuideServer) ListFeatures(rectangle *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.features {
		if utils.InRange(feature.Location, rectangle) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	startTime := time.Now()
	var pointCount, distance int32
	var prevPoint *pb.Point
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:  pointCount,
				Distance:    distance,
				ElapsedTime: int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		if prevPoint != nil {
			distance += utils.CalcDistance(prevPoint, point)
		}
		prevPoint = point
	}
}

func (s *routeGuideServer) recommendOnce(request *pb.RecommendationRequest) (*pb.Feature, error) {
	var nearest, farthest *pb.Feature
	var nearestDistance, farthestDistance int32

	for _, feature := range s.features {
		distance := utils.CalcDistance(feature.Location, request.Point)
		if nearest == nil || distance < nearestDistance {
			nearestDistance = distance
			nearest = feature
		}
		if farthest == nil || distance > farthestDistance {
			farthestDistance = distance
			farthest = feature
		}
	}
	if request.Mode == pb.RecommendationMode_GetFarthest {
		return farthest, nil
	} else {
		return nearest, nil
	}
}

func (s *routeGuideServer) Recommend(stream pb.RouteGuide_RecommendServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		recommended, err := s.recommendOnce(request)
		if err != nil {
			return err
		}
		err = stream.Send(recommended)
		if err != nil {
			return err
		}
		fmt.Println(recommended)
	}
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
