package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/template-service/genproto"
	l "github.com/template-service/pkg/logger"
	cl "github.com/template-service/service/grpc_client"
	"github.com/template-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client: client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.UserResponse, error) {
	user, err := s.storage.User().CreateUser(req)
	if err != nil {
		s.logger.Error(`Filed while inserting user`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while inserting user` )
	}
	postt := req.Posts
	post, err := s.client.PostService().CreatePost(ctx,postt)
	
	user.Posts = append(user.Posts,post)
	return &pb.UserResponse{
		UserId: user.UserId,
		UserName: user.UserName,
		Email:    user.Email,
		Posts: user.Posts, 
	}, nil

}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserName) (*pb.BoolResponse, error) {
	res, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.logger.Error(`Filed while updating user_name`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while updating user_name` )
	}
	return res, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.UserByIdRequest) (*pb.UserResponse, error) {
	user, err := s.storage.User().GetUserById(req.Id)
	if err != nil {
		s.logger.Error(`Filed while getting user by id`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while getting user by id` )
	}
	posts, err := s.client.PostService().GetUserPosts(ctx, &pb.GetUserPostsReq{UserId: user.Id})
	if err != nil {
		s.logger.Error(`Filed while getting posts user by id`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while getting posts user by id` )
	}

	
	return &pb.UserResponse{
		UserId: user.Id,
		UserName: user.UserName,
		Email:    user.Email,
		Posts: posts.Posts, 
	}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.GetAllUsersResponse, error) {
	users, err := s.storage.User().GetAllUsers()
	if err != nil {
		s.logger.Error(`Filed while getting users`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while getting users` )
	}
	return &pb.GetAllUsersResponse{Users: users}, nil
}


