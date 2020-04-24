package service

import (
	"context"

	"github.com/euvsvirus-banan/backend/internal/storage"
	"github.com/euvsvirus-banan/backend/internal/version"
	"github.com/euvsvirus-banan/backend/requests/rpc/requestspb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	logger   *logrus.Entry
	requests *storage.RequestsStorage
}

func New(logger *logrus.Entry, requestData *storage.RequestsStorage) *Service {
	return &Service{
		logger:   logger,
		requests: requestData,
	}
}

func (svc *Service) GetVersion(ctx context.Context, req *requestspb.GetVersionRequest) (*requestspb.GetVersionResponse, error) {
	return &requestspb.GetVersionResponse{
		Project:     version.Project,
		Version:     version.Version,
		BuildDate:   version.BuildDate,
		GitRevision: version.GitRevision,
		GoVersion:   version.GoVersion,
	}, nil
}

func (svc *Service) AddRequest(ctx context.Context, req *requestspb.AddRequestRequest) (*requestspb.AddRequestResponse, error) {
	id := uuid.New().String()
	if err := svc.requests.Add(id, req.Request); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &requestspb.AddRequestResponse{
		RequestId: id,
	}, nil
}

func (svc *Service) DeleteRequest(ctx context.Context, req *requestspb.DeleteRequestRequest) (*requestspb.DeleteRequestResponse, error) {
	if err := svc.requests.Delete(req.RequestId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &requestspb.DeleteRequestResponse{}, nil
}

func (svc *Service) UpdateRequest(ctx context.Context, req *requestspb.UpdateRequestRequest) (*requestspb.UpdateRequestResponse, error) {
	if err := svc.requests.Update(req.RequestId, req.Request); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	u, err := svc.requests.Get(req.RequestId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &requestspb.UpdateRequestResponse{Request: u}, nil
}

func (svc *Service) GetRequests(req *requestspb.GetRequestsRequest, stream requestspb.RequestsRPC_GetRequestsServer) error {
	for id, request := range svc.requests.All() {
		if err := stream.Send(
			&requestspb.GetRequestsResponse{
				RequestId: id,
				Request:   request,
			},
		); err != nil {
			return status.Error(codes.Unknown, err.Error())
		}
	}
	return nil
}
