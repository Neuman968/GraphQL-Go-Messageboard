package services

import (
	"database/sql"
	"fmt"
	"strconv"

	dm "messageboard.example.graphql/.gen/messageboardDB/public/model"

	. "github.com/go-jet/jet/v2/postgres"
	. "messageboard.example.graphql/.gen/messageboardDB/public/table"
)

type PostService struct {
	db *sql.DB
}

func NewPost(db *sql.DB) *PostService {
	return &PostService{
		db: db,
	}
}

func (s *PostService) AddPost(post *dm.Post) (*dm.Post, error) {

	stmt := Post.INSERT(Post.MutableColumns).MODEL(post).
		RETURNING(Post.AllColumns)

	dest := dm.Post{}

	err := stmt.Query(s.db, &dest)

	return &dest, err
}

func (s *PostService) GetPostById(postId string) (*dm.Post, error) {

	postInt, err := strconv.Atoi(postId)

	if err != nil {
		return nil, err
	}

	stmt := Post.SELECT(Post.AllColumns).FROM(Post).
		WHERE(Post.ID.EQ(Int(int64(postInt))))

	dest := dm.Post{}

	err = stmt.Query(s.db, &dest)

	return &dest, err
}

func (s *PostService) AddComment(comment *dm.Comment) (*dm.Comment, error) {
	stmt := Comment.INSERT(Comment.MutableColumns).MODEL(comment).
		RETURNING(Comment.AllColumns)

	dest := dm.Comment{}

	err := stmt.Query(s.db, &dest)

	return &dest, err
}

func (s *PostService) GetPosts() ([]*dm.Post, error) {

	stmt := SELECT(Post.AllColumns).FROM(Post)

	var dest []*dm.Post

	err := stmt.Query(s.db, &dest)

	if err != nil {
		fmt.Printf("Error: %v \n", err)
		return nil, err
	}
	return dest, nil
}
