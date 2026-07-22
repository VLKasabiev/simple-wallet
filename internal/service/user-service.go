package service

import (
	"context"
	"log/slog"

	"github.com/VLKasabiev/simple-wallet/internal/config"
	"github.com/VLKasabiev/simple-wallet/internal/model"
	"github.com/VLKasabiev/simple-wallet/internal/utils"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserService struct {
	repo UserRepository
	cfg  *config.Config
}

func NewUserService(repo UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo: repo,
		cfg: cfg,
	}
}


func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	
    if !utils.CheckPasswordHash(password, user.Password) {
        return "", model.ErrInvalidPassword
    }

	token, err := model.GenerateToken(
		user.ID,
		s.cfg.JWT.SecretKey,
		s.cfg.JWT.ExpiresIn,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*model.User, error) {
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		slog.Warn("user registration failed: email already taken", "email", email)
		return nil, model.ErrEmailAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, err
	}

	user := &model.User{
		Name:  name,
		Email: email,
		Password: hashedPassword,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, IdFromParam, IdFromToken int) (*model.User, error) {
	if IdFromParam != IdFromToken {
        return nil, model.ErrNotUserProfileOwner
    }
	return s.repo.GetByID(ctx, IdFromParam)
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}