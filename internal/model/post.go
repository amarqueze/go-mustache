package model

// Post mirrors the post resource returned by JSONPlaceholder.
type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
