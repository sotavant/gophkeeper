package user

import (
	"context"
	"errors"
	"gophkeeper/internal"
	"gophkeeper/internal/server/auth"
	domain2 "gophkeeper/server/domain"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type ContextUserIDKey struct{}

type Service struct {
	userRepo Repository
}

type Repository interface {
	GetByLogin(ctx context.Context, login string) (domain2.User, error)
	Store(ctx context.Context, user domain2.User) (uint64, error)
}

func NewService(u Repository) *Service {
	return &Service{
		userRepo: u,
	}
}

func (u *Service) Register(ctx context.Context, user domain2.User) (string, error) {
	dbUser, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		internal.Logger.Infow("error in get by login", "err", err)
		return "", domain2.ErrInternalServerError
	}

	if (err != nil && errors.Is(err, pgx.ErrNoRows)) || (err == nil && dbUser.ID != 0) {
		return "", domain2.ErrLoginExist
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		internal.Logger.Infow("error in crypt passwd", "err", err)
		return "", domain2.ErrInternalServerError
	}

	userID, err := u.userRepo.Store(ctx, user)
	if err != nil {
		internal.Logger.Infow("error save user", "err", err)
		return "", domain2.ErrInternalServerError
	}

	token, err := auth.BuildJWTString(userID)
	if err != nil {
		internal.Logger.Infow("error generation token", "err", err)
		return "", domain2.ErrInternalServerError
	}

	return token, nil
}

func (u *Service) Login(ctx context.Context, user domain2.User) (string, error) {
	dbUser, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		internal.Logger.Infow("error in get by login", "err", err)
		return "", domain2.ErrInternalServerError
	}

	if dbUser.ID == 0 {
		return "", domain2.ErrUserNotFound
	}

	passwordCorrect, err := checkPassword(user.Password, dbUser.Password)
	if err != nil {
		internal.Logger.Infow("error in check passwd", "err", err)
		return "", domain2.ErrInternalServerError
	}

	if !passwordCorrect {
		return "", domain2.ErrUserNotFound
	}

	token, err := auth.BuildJWTString(dbUser.ID)
	if err != nil {
		internal.Logger.Infow("error generation token", "err", err)
		return "", domain2.ErrInternalServerError
	}

	return token, nil
}

func HashPassword(password string) (string, error) {
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
