package service

import (
	"context"

	"github.com/euvsvirus-banan/backend/internal/storage"
	"github.com/euvsvirus-banan/backend/internal/version"
	"github.com/euvsvirus-banan/backend/users/rpc/userspb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	logger *logrus.Entry
	users  *storage.UsersStorage
}

func New(logger *logrus.Entry, userData *storage.UsersStorage) *Service {
	return &Service{
		logger: logger,
		users:  userData,
	}
}

func (svc *Service) GetVersion(ctx context.Context, req *userspb.GetVersionRequest) (*userspb.GetVersionResponse, error) {
	return &userspb.GetVersionResponse{
		Project:     version.Project,
		Version:     version.Version,
		BuildDate:   version.BuildDate,
		GitRevision: version.GitRevision,
		GoVersion:   version.GoVersion,
	}, nil
}

func (svc *Service) AddUser(ctx context.Context, req *userspb.AddUserRequest) (*userspb.AddUserResponse, error) {
	id := uuid.New().String()
	if err := svc.users.Add(id, req.User); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userspb.AddUserResponse{
		UserId: id,
	}, nil
}

func (svc *Service) DeleteUser(ctx context.Context, req *userspb.DeleteUserRequest) (*userspb.DeleteUserResponse, error) {
	if err := svc.users.Delete(req.UserId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userspb.DeleteUserResponse{}, nil
}

func (svc *Service) UpdateUser(ctx context.Context, req *userspb.UpdateUserRequest) (*userspb.UpdateUserResponse, error) {
	if err := svc.users.Update(req.UserId, req.User); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	u, err := svc.users.Get(req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userspb.UpdateUserResponse{User: u}, nil
}

func (svc *Service) GetUsers(req *userspb.GetUsersRequest, stream userspb.UsersRPC_GetUsersServer) error {
	for id, user := range svc.users.All() {
		if err := stream.Send(
			&userspb.GetUsersResponse{
				UserId: id,
				User:   user,
			},
		); err != nil {
			return status.Error(codes.Unknown, err.Error())
		}
	}
	return nil
}
