package app

import (
	"time"
)

type Post struct {
	PostId   string    `json:"postId"`
	SharedAt time.Time `json:"sharedAt"`
}
