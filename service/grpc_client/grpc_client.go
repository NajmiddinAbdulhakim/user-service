package grpcClient

import (
	"fmt"

	"github.com/template-service/config"
	pb "github.com/template-service/genproto"
	"google.golang.org/grpc"
)

//GrpcClientI ...
type GrpcClientI interface {
    PostService() pb.PostServiceClient
}

//GrpcClient ...
type GrpcClient struct {
	cfg         config.Config
	postService pb.PostServiceClient
}

//New ...
func New(cfg config.Config) (*GrpcClient, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("post service dial host: %s port: %d",
			cfg.PostServiceHost, cfg.PostServicePort)
	}

	return &GrpcClient{
		cfg:         cfg,
		postService: pb.NewPostServiceClient(connPost),
	}, nil
}

func (s *GrpcClient) PostService() pb.PostServiceClient {
    return s.postService
}
