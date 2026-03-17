package users

import (
	"context"
	repo "kwadw0/gocommerce/internal/adapters/postgres/sqlc"
	"kwadw0/gocommerce/internal/auth"

	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateUser(ctx context.Context, payload CreateUserDTO) (UserResponse, error)
}

type userService struct {
	repo repo.Querier
}

func NewUserService(repo repo.Querier) Service {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, payload CreateUserDTO) (UserResponse, error) {
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
		Phone:     pgtype.Text{String: payload.Phone, Valid: payload.Phone != ""},
	})
	if err != nil {
		return UserResponse{}, err
	}

	return mapToUserResponse(user), nil
}

func mapToUserResponse(user repo.User) UserResponse {
	return UserResponse{
		Uuid:      user.Uuid,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone.String,
	}
}
