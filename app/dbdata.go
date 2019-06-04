package app

import (
	"time"
)

type Share struct {
	UserId   string    `json:"userId"`
	PostId   string    `json:"postId"`
	SharedAt time.Time `json:"sharedAt"`
}

type response struct {
	notifatedUser string `json:"notificated_user"`
	postId        string `json:"post_id"`
}
