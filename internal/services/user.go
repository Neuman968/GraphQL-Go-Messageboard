package services

import (
	"database/sql"

	dm "messageboard.example.graphql/.gen/messageboardDB/public/model"

	. "github.com/go-jet/jet/v2/postgres"
	. "messageboard.example.graphql/.gen/messageboardDB/public/table"
)

type UserService struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetUsers() ([]*dm.Users, error) {
	stmt := SELECT(Users.AllColumns).FROM(Users)

	var dest []*dm.Users
	err := stmt.Query(s.db, &dest)

	if err != nil {
		return nil, err
	}
	return dest, nil
}
