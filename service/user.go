package service

import (
	"context"

	pb "github.com/NajmiddinAbdulhakim/user-service/genproto"
	"github.com/jmoiron/sqlx"
	// user "github.com/NajmiddinAbdulhakim/user-service/genproto"
	l "github.com/NajmiddinAbdulhakim/user-service/pkg/logger"
	cl "github.com/NajmiddinAbdulhakim/user-service/service/grpc_client"
	"github.com/NajmiddinAbdulhakim/user-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().CreateUser(req)
	if err != nil {
		s.logger.Error(`Filed while inserting user`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while inserting user`)
	}
	for _, postt := range req.Posts {
		postt.UserId = user.Id
		post, err := s.client.PostService().CreatePost(ctx, postt)
		if err != nil {
			s.logger.Error(`Filed while inserting user`, l.Error(err))
			return nil, status.Error(codes.Internal, `Filed while inserting user`)
		}
		user.Posts = append(user.Posts, post)
	}

	return user, nil

}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	res, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.logger.Error(`Filed while updating user`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while updating user`)
	}
	return &pb.UpdateUserRes{Success: res}, nil
}

func (s *UserService) GetUserByIdWithPosts(ctx context.Context, req *pb.UserByIdReq) (*pb.User, error) {
	user, err := s.storage.User().GetUserById(req.Id)
	if err != nil {
		s.logger.Error(`Filed while getting user by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting user by id`)
	}
	posts, err := s.client.PostService().GetUserPosts(ctx, &pb.GetUserPostsReq{UserId: user.Id})
	if err != nil {
		s.logger.Error(`Filed while getting posts user by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting posts user by id`)
	}

	user.Posts = posts.Posts

	return user, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *pb.UserByIdReq) (*pb.User, error) {
	user, err := s.storage.User().GetUserById(req.Id)
	if err != nil {
		s.logger.Error(`Filed while getting user by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting user by id`)
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.UserByIdReq) (*pb.UpdateUserRes, error) {
	res, err := s.storage.User().DeleteUser(req.Id)
	if err != nil {
		s.logger.Error(`Filed while delete user`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while delete user`)
	}
	return &pb.UpdateUserRes{Success: res}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *pb.Empty) (*pb.GetAllUsersResponse, error) {
	users, err := s.storage.User().GetAllUsers()
	if err != nil {
		s.logger.Error(`Filed while getting users`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting users`)
	}
	var us []*pb.User
	for _, user := range users {
		posts, err := s.client.PostService().GetUserPosts(ctx, &pb.GetUserPostsReq{UserId: user.Id})
		if err != nil {
			s.logger.Error(`Filed while getting user posts`, l.Error(err))
			return nil, status.Error(codes.Internal, `Filed while getting user posts`)
		}
		user.Posts = posts.Posts
		us = append(us, user)
	}
	return &pb.GetAllUsersResponse{Users: us}, nil
}

func (s *UserService) GetListUsers(ctx context.Context, req *pb.GetUserListReq) (*pb.GetUserListRes, error) {
	users, count, err := s.storage.User().GetListUsers(req.Page, req.Limit)
	if err != nil {
		s.logger.Error(`Filed while getting users list`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting users list`)
	}
	var us []*pb.User
	for _, user := range users {
		posts, err := s.client.PostService().GetUserPosts(ctx, &pb.GetUserPostsReq{UserId: user.Id})
		if err != nil {
			s.logger.Error(`Filed while getting user posts`, l.Error(err))
			return nil, status.Error(codes.Internal, `Filed while getting user posts`)
		}
		user.Posts = posts.Posts
		us = append(us, user)
	}
	return &pb.GetUserListRes{Users: us, Count: count}, nil

}

func (s *UserService) CheckUnique(ctx context.Context, req *pb.CheckUniqueReq) (*pb.CheckUniqueResp, error) {
	exists, err := s.storage.User().CheckUnique(req.Field, req.Value)
	if err != nil {
		s.logger.Error(`Filed check unique for user data`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed to check unique for user data`)
	}
	return &pb.CheckUniqueResp{IsExists: exists}, nil
}