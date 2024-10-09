package gqlmodel

type Comment struct {
	ID           string `json:"id"`
	PostID       string `json:"postId"`
	Post         *Post  `json:"post"`
	AuthorUserID int    `json:"authorUserId"`
	Text         string `json:"text"`
}
