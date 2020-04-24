package service

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/euvsvirus-banan/backend/internal/storage"
	"github.com/euvsvirus-banan/backend/requests/rpc/requestspb"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ws struct {
	b strings.Builder
}

func (w *ws) Seek(offset int64, whence int) (int64, error) {
	w.b = strings.Builder{}
	return 0, nil
}

func (w *ws) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}

func getTestService(w io.WriteSeeker, logger *logrus.Entry) *Service {
	return &Service{
		logger: logger,
		requests: storage.NewRequestsStorage(
			w,
			map[string]*requestspb.Request{
				"a": {
					Title:       "help with groceries",
					Body:        "need help with groceries. I need...",
					RequesterId: "Brown",
					State:       requestspb.Request_WAITING,
				},
				"b": {
					Title:       "help with walking the dog",
					Body:        "need help walking the dog",
					RequesterId: "Brown",
					State:       requestspb.Request_WAITING,
					Skills:      []string{"dog_whisperer"},
					Answers: []*requestspb.Request_Answer{
						{
							VolunteerId: "Blue",
							Comment:     "I love dogs, I have hundreds of 'em!!!",
						},
					},
				},
				"c": {
					Title:       "help with walking the dog",
					Body:        "need help walking the dog",
					RequesterId: "Brown",
					VolunteerId: "Blue",
					State:       requestspb.Request_ACCEPTED,
					Skills:      []string{"dog_whisperer"},
					Answers: []*requestspb.Request_Answer{
						{
							VolunteerId: "Blue",
							Comment:     "I love dogs, I have hundreds of 'em!!!",
						},
					},
				},
				"d": {
					Title:       "help with loneliness",
					Body:        "I fell alone",
					RequesterId: "Brown",
					State:       requestspb.Request_CANCELLED,
					Skills:      []string{"people"},
				},
				"e": {
					Title:       "help with walking the dog",
					Body:        "need help walking the dog",
					RequesterId: "Brown",
					VolunteerId: "Blue",
					State:       requestspb.Request_COMPLETED,
					Skills:      []string{"dog_whisperer"},
					Answers: []*requestspb.Request_Answer{
						{
							VolunteerId: "Blue",
							Comment:     "I love dogs, I have hundreds of 'em!!!",
						},
					},
				},
			},
		),
	}
}

func TestAnswerRequest(t *testing.T) {
	cases := []struct {
		name string
		req  *requestspb.AnswerRequestRequest
		err  error
	}{
		{
			name: "new response",
			req: &requestspb.AnswerRequestRequest{
				RequestId: "a",
				Answer: &requestspb.Request_Answer{
					VolunteerId: "Green",
					Comment:     "I can get you your groceries",
				},
			},
		},
		{
			name: "second response",
			req: &requestspb.AnswerRequestRequest{
				RequestId: "b",
				Answer: &requestspb.Request_Answer{
					VolunteerId: "Green",
					Comment:     "I like dogs even more!",
				},
			},
		},
		{
			name: "already accepted",
			req: &requestspb.AnswerRequestRequest{
				RequestId: "c",
				Answer:    &requestspb.Request_Answer{},
			},
			err: status.Error(codes.InvalidArgument, "request already answered"),
		},
		{
			name: "not found",
			req: &requestspb.AnswerRequestRequest{
				RequestId: "asdasdasd",
				Answer:    &requestspb.Request_Answer{},
			},
			err: status.Error(codes.NotFound, "request not found"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := &ws{}
			logger := logrus.NewEntry(logrus.New())
			svc := getTestService(w, logger)
			_, err := svc.AnswerRequest(context.Background(), tc.req)
			if !cmp.Equal(tc.err, err) {
				t.Error(cmp.Diff(tc.err, err))
			}
		})
	}
}

func TestAcceptHelp(t *testing.T) {
	cases := []struct {
		name string
		req  *requestspb.AcceptHelpRequest
		err  error
	}{
		{
			name: "not found",
			req: &requestspb.AcceptHelpRequest{
				RequestId:   "a",
				VolunteerId: "Green",
			},
			err: status.Error(codes.NotFound, "answer not found"),
		},
		{
			name: "accept",
			req: &requestspb.AcceptHelpRequest{
				RequestId:   "b",
				VolunteerId: "Blue",
			},
		},
		{
			name: "already accepted",
			req: &requestspb.AcceptHelpRequest{
				RequestId: "c",
			},
			err: status.Error(codes.InvalidArgument, "help already accepted or request cancelled"),
		},
		{
			name: "cancelled",
			req: &requestspb.AcceptHelpRequest{
				RequestId: "d",
			},
			err: status.Error(codes.InvalidArgument, "help already accepted or request cancelled"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := &ws{}
			logger := logrus.NewEntry(logrus.New())
			svc := getTestService(w, logger)
			_, err := svc.AcceptHelp(context.Background(), tc.req)
			if !cmp.Equal(tc.err, err) {
				t.Error(cmp.Diff(tc.err, err))
			}
		})
	}
}

func TestCompleteHelp(t *testing.T) {
	cases := []struct {
		name string
		req  *requestspb.CompleteHelpRequest
		err  error
	}{
		{
			name: "not accepted",
			req: &requestspb.CompleteHelpRequest{
				RequestId: "a",
			},
			err: status.Error(codes.InvalidArgument, "help isn't accepted"),
		},
		{
			name: "accepted",
			req: &requestspb.CompleteHelpRequest{
				RequestId: "c",
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := &ws{}
			logger := logrus.NewEntry(logrus.New())
			svc := getTestService(w, logger)
			_, err := svc.CompleteHelp(context.Background(), tc.req)
			if !cmp.Equal(tc.err, err) {
				t.Error(cmp.Diff(tc.err, err))
			}
		})
	}
}

func TestCancelHelp(t *testing.T) {
	cases := []struct {
		name string
		req  *requestspb.CancelHelpRequest
		err  error
	}{
		{
			name: "not accepted",
			req: &requestspb.CancelHelpRequest{
				RequestId: "a",
			},
		},
		{
			name: "completed",
			req: &requestspb.CancelHelpRequest{
				RequestId: "e",
			},
			err: status.Error(codes.InvalidArgument, "request can't be cancelled if it's been completed already"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			w := &ws{}
			logger := logrus.NewEntry(logrus.New())
			svc := getTestService(w, logger)
			_, err := svc.CancelHelp(context.Background(), tc.req)
			if !cmp.Equal(tc.err, err) {
				t.Error(cmp.Diff(tc.err, err))
			}
		})
	}
}
