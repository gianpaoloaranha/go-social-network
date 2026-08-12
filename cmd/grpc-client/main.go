package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/gianpaoloaranha/go-social-network/internal/infra/config"
	socialnetwork "github.com/gianpaoloaranha/go-social-network/proto/gen"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient("localhost:"+cfg.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userClient := socialnetwork.NewUserServiceClient(conn)
	postClient := socialnetwork.NewPostServiceClient(conn)

	user, err := userClient.CreateUser(ctx, &socialnetwork.CreateUserRequest{
		Name:     "Ada Lovelace",
		Email:    "ada@example.com",
		Password: "secret123",
	})
	if err != nil {
		log.Fatalf("create user: %s", status.Convert(err).Message())
	}
	log.Printf("user created: id=%s name=%s email=%s", user.UserId, user.Name, user.Email)

	post, err := postClient.CreatePost(ctx, &socialnetwork.CreatePostRequest{
		AuthorId:    user.UserId,
		Description: "Hello from the grpc client",
	})
	if err != nil {
		log.Fatalf("create post: %s", status.Convert(err).Message())
	}
	log.Printf("post created: id=%s authorId=%s description=%s", post.PostId, post.AuthorId, post.Description)
}
