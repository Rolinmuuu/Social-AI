package service

import "fmt"

var (
	ErrUserAlreadyExisted = fmt.Errorf("user already exists")
	ErrInvalidUser = fmt.Errorf("invalid user or password")
	ErrPostNotFound = fmt.Errorf("post not found")
	ErrCommentNotFound = fmt.Errorf("comment not found")
	ErrLikeNotFound = fmt.Errorf("like not found")
	ErrShareNotFound = fmt.Errorf("share not found")
	ErrFollowNotFound = fmt.Errorf("follow not found")
	ErrMessageNotFound = fmt.Errorf("message not found")
	ErrUserNotFound = fmt.Errorf("user not found")
)