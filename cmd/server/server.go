package server

import (
	"context"
	"log"
	"net"

	pb "github.com/hi20160616/fetchnews-api/proto/v1"
	"github.com/hi20160616/ms-bbc/config"
	"github.com/hi20160616/ms-bbc/internal/fetcher"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", config.Data.MS.Addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFetchNewsServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

type server struct {
	pb.UnimplementedFetchNewsServer
}

func (s *server) List(ctx context.Context, in *pb.ListArticlesRequest) (*pb.ListArticlesResponse, error) {
	log.Printf("Received: %v", in.GetPageSize())
	a := fetcher.NewArticle()
	as, err := a.List()
	if err != nil {
		return nil, err
	}
	resp := &pb.ListArticlesResponse{}
	for _, a := range as {
		resp.Articles = append(resp.Articles, &pb.Article{
			Id:            a.Id,
			Title:         a.Title,
			Content:       a.Content,
			WebsiteId:     a.WebsiteId,
			WebsiteTitle:  a.WebsiteTitle,
			WebsiteDomain: a.WebsiteDomain,
			UpdateTime:    a.UpdateTime,
		})
	}
	return resp, nil
}

func (s *server) Get(ctx context.Context, in *pb.GetArticleRequest) (*pb.Article, error) {
	log.Printf("Id: %v", in.Id)
	// Got article via json reading
	a := fetcher.NewArticle()
	a, err := a.Get(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Article{
		Id:            a.Id,
		Title:         a.Title,
		Content:       a.Content,
		WebsiteId:     a.WebsiteId,
		WebsiteTitle:  a.WebsiteTitle,
		WebsiteDomain: a.WebsiteDomain,
		UpdateTime:    a.UpdateTime,
	}, nil
}
