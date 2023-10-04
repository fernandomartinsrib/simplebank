package gapi

import (
	"fmt"

	db "github.com/fernandomartinsrib/simplebank/db/sqlc"
	"github.com/fernandomartinsrib/simplebank/pb"
	"github.com/fernandomartinsrib/simplebank/token"
	"github.com/fernandomartinsrib/simplebank/utils"
	"github.com/fernandomartinsrib/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          utils.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
