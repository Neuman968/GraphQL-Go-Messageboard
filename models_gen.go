// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package main

type AddNewCommentInput struct {
	PostID string `json:"postId"`
	Text   string `json:"text"`
}

type AddNewPostInput struct {
	Text string `json:"text"`
}

type Comment struct {
	ID           string `json:"id"`
	PostID       string `json:"postId"`
	Post         *Post  `json:"post"`
	AuthorUserID int    `json:"authorUserId"`
	AuthorUser   *User  `json:"authorUser"`
	Text         string `json:"text"`
}

type Mutation struct {
}

type Post struct {
	ID           string     `json:"id"`
	AuthorUserID string     `json:"authorUserId"`
	AuthorUser   *User      `json:"authorUser"`
	Text         string     `json:"text"`
	Comments     []*Comment `json:"comments"`
}

type Query struct {
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
