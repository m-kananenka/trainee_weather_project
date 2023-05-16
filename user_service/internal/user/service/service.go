package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
	"user_service/internal/config"
	"user_service/internal/user/model"
)

type repository interface {
	AddUser(ctx context.Context, modelUser model.User) error
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

func (c *Controller) Create(context context.Context, user *model.User) (*model.User, error) {
	if user.Name == "" {
		return nil, errors.New("name is a vital field")
	}
	user.ID = uuid.New().String()
	i := *user
	return nil, c.repo.AddUser(context, i)
}

func (c *Controller) GetUser(context context.Context, id string) (model.User, error) {
	return c.repo.GetUser(context, id)
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
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
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
