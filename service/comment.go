package service

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"

	"socialai/constants"
	"socialai/backend"
	"socialai/model"
)

func AddComment(postId, parentCommentId, userId, content string) (string, error) {
	if postId == "" || userId == "" || strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("Invalid comment input")
	}
	
	postQuery := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("post_id", postId)).
		MustNot(elastic.NewTermQuery("deleted", true))
	postSearchResult, err := backend.ESBackend.ReadFromES(postQuery, constants.POST_INDEX)
	if err != nil {
		return "", err
	}
	posts := getPostFromSearchResult(postSearchResult)
	if len(posts) == 0 {
		return "", fmt.Errorf("Post not found")
	}

	commentId := uuid.New()
	now := time.Now().Unix()

	rootCommentId := commentId
	commentLevel := 0

	if parentCommentId != "" {
		parentCommentQuery := elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("comment_id", parentCommentId)).
			MustNot(elastic.NewTermQuery("deleted", true))
		parentCommentSearchResult, err := backend.ESBackend.ReadFromES(parentCommentQuery, constants.COMMENT_INDEX)
		if err != nil {
			return "", err
		}
		parentComments := getCommentFromSearchResult(parentCommentSearchResult)
		if len(parentComments) == 0 {
			return "", fmt.Errorf("Parent comment not found")
		}
		parentComment := parentComments[0]

		if parentComment.PostId != postId { 
			return "", fmt.Errorf("Parent comment does not belong to the post")
		}

		rootCommentId = parentComment.RootCommentId
		if rootCommentId == "" {
			rootCommentId = parentComment.CommentId
		}
		depth = parentComment.Depth + 1
	}



	comment := model.Comment{
		CommentId: commentId,
		ParentCommentId: parentCommentId,
		RootCommentId: rootCommentId,
		UserId: userId,
		PostId: postId,
		Depth: depth,
		Content: content,
		CreatedAt: now,
		Deleted: false,
		DeletedAt: 0,
	}
	if err := backend.ESBackend.SaveToES(comment, constants.COMMENT_INDEX, comment.CommentId); err != nil {
		return "", err
	}

	return commentId, nil
}

func getCommentFromSearchResult(searchResult *elastic.SearchResult) []model.Comment {
	var comments []model.Comment
	var ctype model.Comment
	for _, item := range searchResult.Each(reflect.TypeOf(ctype)) {
		if c, ok := item.(model.Comment); ok && !c.Deleted {
			comments = append(comments, c)
		}
	}
	return comments
}