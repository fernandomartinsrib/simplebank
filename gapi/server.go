package gapi

import (
	"fmt"

	db "github.com/fernandomartinsrib/simplebank/db/sqlc"
	"github.com/fernandomartinsrib/simplebank/pb"
	"github.com/fernandomartinsrib/simplebank/token"
	"github.com/fernandomartinsrib/simplebank/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
