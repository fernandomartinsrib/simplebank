package gapi

import (
	"context"

	db "github.com/fernandomartinsrib/simplebank/db/sqlc"
	"github.com/fernandomartinsrib/simplebank/pb"
	"github.com/fernandomartinsrib/simplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify user's email")
	}

	resp := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	return resp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := validation.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
