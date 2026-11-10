package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gomustache/internal/model"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ErrNotFound is returned when the remote API answers 404.
var ErrNotFound = errors.New("resource not found")

type JSONPlaceholderClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewJSONPlaceholderClient(baseURL string, httpClient *http.Client) *JSONPlaceholderClient {
	return &JSONPlaceholderClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpClient,
	}
}

func (c *JSONPlaceholderClient) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := c.getJSON(ctx, "/users", &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (c *JSONPlaceholderClient) GetUser(ctx context.Context, id int) (*model.User, error) {
	var u model.User
	if err := c.getJSON(ctx, fmt.Sprintf("/users/%d", id), &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *JSONPlaceholderClient) ListPostsByUser(ctx context.Context, userID int) ([]model.Post, error) {
	q := url.Values{}
	q.Set("userId", strconv.Itoa(userID))

	var posts []model.Post
	if err := c.getJSON(ctx, "/posts?"+q.Encode(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// getJSON performs a GET request and decodes the JSON body into out.
func (c *JSONPlaceholderClient) getJSON(ctx context.Context, path string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("build request %s: %w", path, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("call %s: unexpected status %d", path, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}
