package server

import (
	"context"
	"log"

	"github.com/JUSSIAR/Golang-gRPC-blog/service/api"
	"github.com/JUSSIAR/Golang-gRPC-blog/service/internal/cache"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type BlogServer struct {
	api.BlogServerServer
	Logger        *zap.Logger
	SimpleCounter *prometheus.CounterVec
	Histogram     *prometheus.HistogramVec
	RedisClent    *redis.Client
	PostgresDB    *gorm.DB
}

func (s *BlogServer) Like(
	ctx context.Context, req *api.LikeRequest,
) (*api.LikeResponse, error) {
	cache.Like(ctx, *s.RedisClent, string(req.EType), req.EId, req.AuthorLogin)
	return &api.LikeResponse{
		LikeCount: 1,
		IsLiked:   true,
	}, nil
}

func (s *BlogServer) Unlike(
	ctx context.Context, req *api.LikeRequest,
) (*api.LikeResponse, error) {
	cache.Unlike(ctx, *s.RedisClent, string(req.EType), req.EId, req.AuthorLogin)
	return &api.LikeResponse{
		LikeCount: 0,
		IsLiked:   false,
	}, nil
}

func (s *BlogServer) CreateComment(
	ctx context.Context, req *api.CommentRequestCreate,
) (*api.CommentResponse, error) {
	log.Println("create post", req)
	return &api.CommentResponse{
		Id:          1,
		AuthorLogin: "admin",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
		PostId:      0,
	}, nil
}

func (s *BlogServer) ReadComment(
	ctx context.Context, req *api.CommentRequestReadDelete,
) (*api.CommentResponse, error) {
	log.Println("read post", req)
	return &api.CommentResponse{
		Id:          1,
		AuthorLogin: "admin",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
		PostId:      0,
	}, nil
}

func (s *BlogServer) UpdateComment(
	ctx context.Context, req *api.CommentRequestUpdate,
) (*api.CommentResponse, error) {
	log.Println("update post", req)
	return &api.CommentResponse{
		Id:          1,
		AuthorLogin: "admin",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
		PostId:      0,
	}, nil
}

func (s *BlogServer) DeleteComment(
	ctx context.Context, req *api.CommentRequestReadDelete,
) (*api.CommentResponse, error) {
	log.Println("delete post", req)
	return &api.CommentResponse{
		Id:          1,
		AuthorLogin: "admin",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
		PostId:      0,
	}, nil
}

func (s *BlogServer) CreatePost(
	ctx context.Context, req *api.PostRequestCreate,
) (*api.PostResponse, error) {
	log.Println("create post", req)
	return &api.PostResponse{
		Id:          0,
		AuthorLogin: "admin",
		Title:       "postTitle",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
	}, nil
}

func (s *BlogServer) ReadPost(
	ctx context.Context, req *api.PostRequestReadDelete,
) (*api.PostResponse, error) {
	log.Println("read post", req)
	return &api.PostResponse{
		Id:          0,
		AuthorLogin: "admin",
		Title:       "postTitle",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
	}, nil
}

func (s *BlogServer) UpdatePost(
	ctx context.Context, req *api.PostRequestUpdate,
) (*api.PostResponse, error) {
	log.Println("update post", req)
	return &api.PostResponse{
		Id:          0,
		AuthorLogin: "admin",
		Title:       "postTitle",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
	}, nil
}

func (s *BlogServer) DeletePost(
	ctx context.Context, req *api.PostRequestReadDelete,
) (*api.PostResponse, error) {
	log.Println("delete post", req)
	return &api.PostResponse{
		Id:          0,
		AuthorLogin: "admin",
		Title:       "postTitle",
		Text:        "postText",
		Timestamp:   timestamppb.Now(),
	}, nil
}
