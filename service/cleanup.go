package service

import (
	"socialai/backend"
	"socialai/constants"

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
	posts := getPostFromSearchResult(searchResult)
	if len(posts) == 0 {
		return false, nil
	}
	if limit <= 0 || limit > len(posts) {
		limit = len(posts)
	}
	for _, post := range posts[:limit] {
		err := backend.GCSBackend.DeleteFromGCS(post.Id)
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
		if saveErr := backend.ESBackend.SaveToES(*post, constants.POST_INDEX, post.Id); saveErr != nil {
			return false, saveErr
		}
	}
	return true, nil
}