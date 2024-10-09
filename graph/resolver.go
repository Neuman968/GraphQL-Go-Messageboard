package graph

import (
	"messageboard.example.graphql/internal/services"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	PostService    *services.PostService
	UserService    *services.UserService
	LoggedInUserId int
}
