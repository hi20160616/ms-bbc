package server

import (
	"context"
	"log"
	"net"

	pb "github.com/hi20160616/fetchnews-api/proto/v1"
	"github.com/hi20160616/ms-bbc/config"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFetchNewsServer
}

func (s *server) List(ctx context.Context, in *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error) {
	log.Printf("Received: %v", in.GetPageSize())
	return &pb.ListArticlesResponse{Articles: nil}, nil
}

func (s *server) Get(ctx context.Context, in *pb.GetArticleRequest) (*pb.Article, error) {
	log.Printf("Id: %v", in.Id)
	a := &pb.Article{Id: in.Id} // Got article via json reading
	return a, nil
}

func main() {
	lis, err := net.Listen("tcp", config.Data.MS["addr"])
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFetchNewsServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}
