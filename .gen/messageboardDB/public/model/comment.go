//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Comment struct {
	ID            int32 `sql:"primary_key"`
	PostID        int32
	AuthorUsersID int32
	CreatedTime   *time.Time
	Text          *string
}
