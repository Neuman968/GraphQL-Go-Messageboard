package gqlmodel

type Post struct {
	ID           string `json:"id"`
	AuthorUserID string `json:"authorUserId"`
	Text         string `json:"text"`
	// Comments     []*Comment `json:"comments"`
}
