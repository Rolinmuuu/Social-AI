package service

import (
	"time"
	"reflect"
	"socialai/backend"
	"socialai/constants"
	"socialai/model"
	"mime/multipart"
	"github.com/olivere/elastic/v7"
)

func SearchPostByUser(user string) ([]model.Post, error) {
	// 1. define business logic 
	baseQuery := elastic.NewTermQuery("user", user)
	query := elastic.NewBoolQuery().
		Must(baseQuery).
		MustNot(elastic.NewTermQuery("deleted", true))
	// 2. call backend
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	// 3. construct response: researchResult -> []model.Post
	return getPostFromSearchResult(searchResult), nil
}

func SearchPostByKeywords(keywords string) ([]model.Post, error) {
	baseQuery := elastic.NewMatchQuery("message", keywords)
	baseQuery.Operator("AND")
	if keywords == "" {
		baseQuery.ZeroTermsQuery("all")
	}
	query := elastic.NewBoolQuery().
		Must(baseQuery).
		MustNot(elastic.NewTermQuery("deleted", true))
	// Set operator to "AND" to make sure all keywords are included in search results.
	
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}

	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var posts []model.Post
	var ptype model.Post
	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		if p, ok := item.(model.Post); ok {
			if p.Deleted {
				continue
			}
			posts = append(posts, p)
		}
	}
	return posts
}

func SavePost(post *model.Post, file multipart.File) error {
	// 1. save the GCS and get URL
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}

	post.Url = medialink
	post.Deleted = false
	post.DeletedAt = 0
	post.CleanupStatus = ""
	post.RetryCount = 0
	post.LastError = ""

	// 2. save the post to database (Elasticsearch)
	// 3. response
	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
	
	// 如果两个并非同时成功，无法保证 Atomic 的要求，怎么办呢
	// 
}

func DeletePost(postID string) (bool, error) {
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("id", postID)).
		MustNot(elastic.NewTermQuery("deleted", true))
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return false, err
	}
	posts := getPostFromSearchResult(searchResult)
	if len(posts) == 0 {
		return false, nil
	}
	post := posts[0]
	post.Deleted = true
	post.DeletedAt = time.Now().Unix()
	post.CleanupStatus = "pending"
	post.RetryCount = 0
	post.LastError = ""
	
	err = backend.ESBackend.SaveToES(&post, constants.POST_INDEX, post.Id)
	if err != nil {
		return false, err
	}
	return true, nil
}