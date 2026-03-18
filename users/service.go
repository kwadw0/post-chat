package users

import (
	"context"
	"errors"
	repo "kwadw0/gocommerce/internal/adapters/postgres/sqlc"
	"kwadw0/gocommerce/internal/auth"
	"kwadw0/gocommerce/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service interface {
	CreateUser(ctx context.Context, payload CreateUserDTO) (UserResponse, error)
	Login(ctx context.Context, payload LoginDto) (UserResponse, string, error)
}

type userService struct {
	repo       repo.Querier
	jwt_secret []byte
	jwt_exp    time.Duration
}

func NewUserService(repo repo.Querier, secret []byte, exp time.Duration) Service {
	return &userService{
		repo:       repo,
		jwt_secret: secret,
		jwt_exp:    exp,
	}
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 is unique_violation
				return UserResponse{}, utils.ErrDuplicateEmail
			}
		}
		return UserResponse{}, err
	}

	return mapToUserResponse(user), nil
}

func (s *userService) Login(ctx context.Context, payload LoginDto) (UserResponse, string, error) {
	exitingUser, err := s.repo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserResponse{}, "", utils.ErrUserNotFound
		}
		return UserResponse{}, "", err
	}

	if err := auth.VerifyPassword(payload.Password, exitingUser.Password); err != nil {
		return UserResponse{}, "", utils.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(s.jwt_secret, exitingUser.Uuid, s.jwt_exp)
	if err != nil {
		return UserResponse{}, "", err
	}

	return mapToUserResponse(exitingUser), token, nil
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
