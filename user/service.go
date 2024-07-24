package user

import (
	"context"
	"errors"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/server/auth"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type ContextUserIDKey struct{}

type Service struct {
	userRepo UserRepository
}

type UserRepository interface {
	GetByLogin(ctx context.Context, login string) (domain.User, error)
	Store(ctx context.Context, user domain.User) (int64, error)
}

func NewService(u UserRepository) *Service {
	return &Service{
		userRepo: u,
	}
}

func (u *Service) Register(ctx context.Context, user domain.User) (string, error) {
	dbUser, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		internal.Logger.Infow("error in get by login", "err", err)
		return "", domain.ErrInternalServerError
	}

	if (err != nil && errors.Is(err, pgx.ErrNoRows)) || (err == nil && dbUser.ID != 0) {
		return "", domain.ErrLoginExist
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		internal.Logger.Infow("error in crypt passwd", "err", err)
		return "", domain.ErrInternalServerError
	}

	userID, err := u.userRepo.Store(ctx, user)
	if err != nil {
		internal.Logger.Infow("error save user", "err", err)
		return "", domain.ErrInternalServerError
	}

	token, err := auth.BuildJWTString(userID)
	if err != nil {
		internal.Logger.Infow("error generation token", "err", err)
		return "", domain.ErrInternalServerError
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func checkPassword(password, passwordHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, err
}
