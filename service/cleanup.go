package service

import (
	"reflect"
	"socialai/backend"
	"socialai/constants"
	"socialai/model"

	"github.com/olivere/elastic/v7"
)

func CleanupDeletedPost(limit int	) (bool, error) {
	query := elastic.NewBoolQuery().
		Must(
			elastic.NewTermQuery("deleted", true),
			elastic.NewTermQuery("cleanup_status", "pending"),
		)
		
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return false, err
	}
	posts := getDeletedPostFromSearchResult(searchResult)
	if len(posts) == 0 {
		return false, nil
	}
	if limit <= 0 || limit > len(posts) {
		limit = len(posts)
	}
	for i := 0; i < limit; i++ {
		post := posts[i]
		err := backend.GCSBackend.DeleteFromGCS(post.PostId)
		if err != nil {
			post.RetryCount++
			post.LastError = err.Error()
			if post.RetryCount >= 5 {
				post.CleanupStatus = "failed"
			} else {
				post.CleanupStatus = "pending"
			} 
		 } else {
			post.CleanupStatus = "completed"
			post.LastError = ""
		}
		if saveErr := backend.ESBackend.SaveToES(&post, constants.POST_INDEX, post.PostId); saveErr != nil {
			return false, saveErr
		}
	}
	return true, nil
}

func getDeletedPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var posts []model.Post
	var ptype model.Post
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		if p, ok := item.(model.Post); ok {
			posts = append(posts, p)
		}
	}
	return posts
}