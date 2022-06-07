package repo

import (
    pb "github.com/template-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
    CreateUser(*pb.User) (*pb.User, error)
    UpdateUser(*pb.UpdateUserReq) (*pb.UpdateUserRes, error)
    GetUserById(userID string) (*pb.User, error)
    GetAllUsers() ([]*pb.User, error)
}
