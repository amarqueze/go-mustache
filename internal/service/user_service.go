package service

import (
	"context"
	"errors"
	"time"

	"gomustache/internal/client"
	"gomustache/internal/model"

	"golang.org/x/sync/errgroup"
)

var ErrUserNotFound = errors.New("user not found")

type UserClient interface {
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUser(ctx context.Context, id int) (*model.User, error)
	ListPostsByUser(ctx context.Context, userID int) ([]model.Post, error)
}

var _ UserClient = (*client.JSONPlaceholderClient)(nil)

type UserService struct {
	client     UserClient
	writeDelay time.Duration
}

func NewUserService(c UserClient, writeDelay time.Duration) *UserService {
	return &UserService{client: c, writeDelay: writeDelay}
}

func (s *UserService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.client.ListUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*model.User, error) {
	u, err := s.client.GetUser(ctx, id)
	if errors.Is(err, client.ErrNotFound) {
		return nil, ErrUserNotFound
	}
	return u, err
}

func (s *UserService) GetUserWithPosts(ctx context.Context, id int) (*model.UserWithPosts, error) {
	g, ctx := errgroup.WithContext(ctx)

	var user *model.User
	var posts []model.Post

	g.Go(func() error {
		var err error
		user, err = s.GetUser(ctx, id)
		return err
	})

	g.Go(func() error {
		var err error
		posts, err = s.client.ListPostsByUser(ctx, id)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &model.UserWithPosts{User: *user, Posts: posts}, nil
}

func (s *UserService) CreateUser(ctx context.Context, in model.UserInput) error {
	return s.simulateWork(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, in model.UserInput) error {
	if _, err := s.GetUser(ctx, id); err != nil {
		return err
	}
	return s.simulateWork(ctx)
}

func (s *UserService) simulateWork(ctx context.Context) error {
	select {
	case <-time.After(s.writeDelay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
