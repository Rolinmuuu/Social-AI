package service

import (
	"reflect"
	"socialai/backend"
	"socialai/constants"
	"socialai/model"
	"mime/multipart"
	"github.com/olivere/elastic/v7"
)

func SearchPostByUser(user string) ([]model.Post, error) {
	// 1. define business logic 
	query := elastic.NewTermQuery("user", user)

	// 2. call backend
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	// 3. construct response: researchResult -> []model.Post
	return getPostFromSearchResult(searchResult), nil
}

func SearchPostByKeywords(keywords string) ([]model.Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	// Set operator to "AND" to make sure all keywords are included in search results.
	// If keywords is empty, it will be treated as a match all query, which returns all posts in database.
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	
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

	// 2. save the post to database (Elasticsearch)
	// 3. response
	return backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
	
	// 如果两个并非同时成功，无法保证 Atomic 的要求，怎么办呢
	// 
}