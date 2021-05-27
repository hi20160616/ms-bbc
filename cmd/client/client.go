package client

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/hi20160616/fetchnews-api/proto/v1"
	"github.com/hi20160616/ms-bbc/config"
	"google.golang.org/grpc"
)

var address = "localhost" + config.Data.MS["addr"]

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFetchNewsClient(conn)

	// Contact the server and print out its response.
	name := "bbc_server"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ListArticles(ctx, &pb.ListArticlesRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetArticles())
	// r, err = c.GetArticle(ctx, &pb.GetArticleRequest{Id: name})
	article, err := c.GetArticle(ctx, &pb.GetArticleRequest{Id: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", article.Title)
}
