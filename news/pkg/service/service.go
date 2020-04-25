package service

import (
	"context"

	"github.com/euvsvirus-banan/backend/internal/storage"
	"github.com/euvsvirus-banan/backend/internal/version"
	"github.com/euvsvirus-banan/backend/news/rpc/newspb"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	logger *logrus.Entry
	news   *storage.NewsStorage
}

func New(logger *logrus.Entry, newData *storage.NewsStorage) *Service {
	return &Service{
		logger: logger,
		news:   newData,
	}
}

func (svc *Service) GetVersion(ctx context.Context, req *newspb.GetVersionRequest) (*newspb.GetVersionResponse, error) {
	return &newspb.GetVersionResponse{
		Project:     version.Project,
		Version:     version.Version,
		BuildDate:   version.BuildDate,
		GitRevision: version.GitRevision,
		GoVersion:   version.GoVersion,
	}, nil
}

func (svc *Service) AddNew(ctx context.Context, req *newspb.AddNewRequest) (*newspb.AddNewResponse, error) {
	id := uuid.New().String()
	if err := svc.news.Add(id, req.New); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &newspb.AddNewResponse{
		NewId: id,
	}, nil
}

func (svc *Service) DeleteNew(ctx context.Context, req *newspb.DeleteNewRequest) (*newspb.DeleteNewResponse, error) {
	if err := svc.news.Delete(req.NewId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &newspb.DeleteNewResponse{}, nil
}

func (svc *Service) UpdateNew(ctx context.Context, req *newspb.UpdateNewRequest) (*newspb.UpdateNewResponse, error) {
	if err := svc.news.Update(req.NewId, req.New); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	u, err := svc.news.Get(req.NewId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &newspb.UpdateNewResponse{New: u}, nil
}

func (svc *Service) GetNews(req *newspb.GetNewsRequest, stream newspb.NewsRPC_GetNewsServer) error {
	for id, new := range svc.news.All() {
		if err := stream.Send(
			&newspb.GetNewsResponse{
				NewId: id,
				New:   new,
			},
		); err != nil {
			return status.Error(codes.Unknown, err.Error())
		}
	}
	return nil
}

func (svc *Service) GetNewsByID(ctx context.Context, req *newspb.GetNewsByIDRequest) (*newspb.GetNewsByIDResponse, error) {
	n, err := svc.news.Get(req.NewsId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "new not found")
	}
	return &newspb.GetNewsByIDResponse{
		News: n,
	}, nil
}

func (svc *Service) SearchNewsByPostcode(req *newspb.SearchNewsByPostcodeRequest, stream newspb.NewsRPC_SearchNewsByPostcodeServer) error {
	for id, n := range svc.news.All() {
		if n.Postcode != req.Postcode {
			continue
		}
		if err := stream.Send(
			&newspb.SearchNewsByPostcodeResponse{
				NewsId: id,
				News:   n,
			},
		); err != nil {
			return status.Error(codes.Unknown, err.Error())
		}
	}
	return nil
}
