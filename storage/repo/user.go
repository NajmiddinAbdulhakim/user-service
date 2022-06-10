package repo

import (
    pb "github.com/template-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
    CreateUser(*pb.User) (*pb.User, error)
    UpdateUser(*pb.UpdateUserReq) (bool, error)
    GetUserById(userID string) (*pb.User, error)
    GetAllUsers() ([]*pb.User, error)
    DeleteUser(userID string)(bool, error)

    GetListUsers(page, limit int64) ([]*pb.User,int64, error)
}
