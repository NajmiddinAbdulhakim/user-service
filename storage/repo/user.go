package repo

import (
    pb "github.com/template-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
    CreateUser(*pb.User) (*pb.UserResponse, error)
    UpdateUser(*pb.UpdateUserName) (*pb.BoolResponse, error)
    GetUserById(userID string) (*pb.User, error)
    GetAllUsers() ([]*pb.User, error)
}
