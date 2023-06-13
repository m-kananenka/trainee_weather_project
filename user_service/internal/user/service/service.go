package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"user_service/api/pb"
	"user_service/internal/config"
	"user_service/internal/user/model"
)

//go:generate mockgen -source ./service.go -destination ../mock/service.go -package mock

type repository interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	CheckAuth(ctx context.Context, login, password string) (model.User, error)
}

type Controller struct {
	repo repository
	cfg  *config.Config
}

func NewController(repo repository, cfg *config.Config) *Controller {
	return &Controller{
		repo: repo,
		cfg:  cfg,
	}
}

func (c *Controller) Create(ctx context.Context, user *model.User) error {
	if user.Name == "" {
		return errors.New("name is a vital field")
	}
	user.ID = uuid.New().String()
	return c.repo.CreateUser(ctx, *user)
}

func (c *Controller) GetUser(ctx context.Context, id string) (model.User, error) {
	return c.repo.GetUser(ctx, id)
}

func (c *Controller) UpdateUser(ctx context.Context, user model.User) error {
	_, err := c.repo.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}
	return c.repo.UpdateUser(ctx, user)
}

func (c *Controller) DeleteUser(ctx context.Context, id string) error {
	_, err := c.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}
	return c.repo.DeleteUser(ctx, id)
}

func (c *Controller) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return c.repo.GetAllUsers(ctx)
}

func (c *Controller) Authorize(ctx context.Context, login, password string) (string, error) {
	user, err := c.repo.CheckAuth(ctx, login, password)
	if err != nil {
		return "", fmt.Errorf("failed to authorize user: %w", err)
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    user.ID,
		Subject:   "authorized",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.cfg.JWTKeyword))
	if err != nil {
		return "", fmt.Errorf("failed to sigh token: %w", err)
	}
	return tokenString, nil
}

//gPRC

func (c *Controller) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	_, err := c.repo.CheckAuth(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		return &pb.Response{Response: false}, err
	}
	return &pb.Response{Response: true}, nil
}
