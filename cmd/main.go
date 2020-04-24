package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/euvsvirus-banan/backend/internal/storage"
	"github.com/euvsvirus-banan/backend/users/pkg/service"
	"github.com/euvsvirus-banan/backend/users/rpc/userspb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getLogger(debug bool) *logrus.Entry { // nolint: unparam
	l := logrus.New()
	if debug {
		l.SetLevel(logrus.DebugLevel)
	} else {
		l.SetLevel(logrus.InfoLevel)
	}
	return logrus.NewEntry(l)
}

func startService(logger *logrus.Entry, addr string, userData *storage.UserStorage) error {
	logger.WithFields(
		logrus.Fields{
			"addr": addr,
		},
	).Info("starting server")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("%w: problem listening on given address", err)
	}

	grpc_logrus.ReplaceGrpcLogger(logger)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logger),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logger),
		),
	)

	svc := service.New(logger, userData)

	userspb.RegisterUsersRPCServer(grpcServer, svc)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("%w: problem serving service", err)
	}
	return nil
}

func getUserData(file io.ReadWriteSeeker) (*storage.UserStorage, error) {
	data := make(map[string]*userspb.User)
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("problem trying to read user data file: %w", err)
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("problem unmarshalling user data: %w", err)
	}
	st := storage.NewUserStorage(file, data)
	return st, nil
}

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	addr := flag.String("addr", "127.0.0.1:65010", "Address to bind the service to")
	usersFile := flag.String("users-file", "/euvsvirus-backend/users.json", "File to store user information")

	flag.Parse()

	logger := getLogger(*debug)

	file, err := os.OpenFile(*usersFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	userData, err := getUserData(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := startService(
		logger,
		*addr,
		userData,
	); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
