package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/JUSSIAR/Golang-gRPC-blog/service/internal/cache"
	"github.com/JUSSIAR/Golang-gRPC-blog/service/storages"
	"github.com/go-redis/redis/v8"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	api "github.com/JUSSIAR/Golang-gRPC-blog/service/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type BlogServer struct {
	api.BlogServerServer
	Logger        *zap.Logger
	SimpleCounter *prometheus.CounterVec
	Histogram     *prometheus.HistogramVec
	redisClent    *redis.Client
}

func (s *BlogServer) Like(
	ctx context.Context, req *api.LikeRequest,
) (*api.LikeResponse, error) {
	cache.Like(ctx, *s.redisClent, string(req.EType), req.EId, req.AuthorLogin)
	return &api.LikeResponse{
		LikeCount: 1,
		IsLiked:   true,
	}, nil
}

func (s *BlogServer) Unlike(
	ctx context.Context, req *api.LikeRequest,
) (*api.LikeResponse, error) {
	cache.Unlike(ctx, *s.redisClent, string(req.EType), req.EId, req.AuthorLogin)
	return &api.LikeResponse{
		LikeCount: 0,
		IsLiked:   false,
	}, nil
}

//go:embed api/api.swagger.json
var swaggerData []byte

//go:embed swagger-ui
var swaggerFiles embed.FS

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "simple_counter",
		},
		[]string{"label"},
	)
	prometheus.MustRegister(counter)

	histCounter := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "histogram_counter",
			Buckets: []float64{0, 10.0, 20.0, 30.0},
		},
		[]string{"label"},
	)
	prometheus.MustRegister(histCounter)

	go func() {
		server := &http.Server{
			Addr:    ":9000",
			Handler: promhttp.Handler(),
		}

		log.Println("Serving metrics on http://0.0.0.0:9000")
		log.Fatalln(server.ListenAndServe())
	}()

	ctx := context.Background()
	postgresDB := storages.ConnectPostgres()
	redisDB := storages.ConnectRedis(ctx)

	fmt.Println("OK")
	fmt.Println(postgresDB)
	fmt.Println(redisDB)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := &BlogServer{
		Logger:        logger,
		SimpleCounter: counter,
		Histogram:     histCounter,
		redisClent:    &redisDB,
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_prometheus.UnaryServerInterceptor,
		),
	)
	api.RegisterBlogServerServer(grpcServer, s)
	reflection.Register(grpcServer)

	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	if err := api.RegisterBlogServerHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(swaggerData)
	})

	fSys, err := fs.Sub(swaggerFiles, "swagger-ui")
	if err != nil {
		panic(err)
	}

	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.FS(fSys))))

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
