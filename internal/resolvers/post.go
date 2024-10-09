package resolvers

import (
	"messageboard.example.graphql/graph/model"
)

type PostResolver struct {
	posts []model.Post
}
