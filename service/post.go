package service

import (
	"fmt"
	"time"
	"reflect"

	"socialai/backend"
	"socialai/constants"
	"socialai/model"
	"socialai/utils"

	"mime/multipart"
	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"
)

func SearchPostByUserId(userId string) ([]model.Post, error) {
	// 1. check cache, and if hit, return the cached posts
	ctx := context.Background()
	cacheKey := utils.userFeedCacheKey(user)
	cacheValue, err := backend.RedisBackend.Get(ctx, cacheKey)
	if err == nil {
		var posts []model.Post
		if err := json.Unmarshal([]byte(cacheValue), &posts); err = nil {
			fmt.Println("cache hit for user feed: %s", userId)
			return posts, nil
		}
		return posts, nil
	}

	// 2. if not hit, define business logic 
	baseQuery := elastic.NewTermQuery("user_id", userId)
	query := elastic.NewBoolQuery().
		Must(baseQuery).
		MustNot(elastic.NewTermQuery("deleted", true))

	// 3. call backend
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	posts := getPostFromSearchResult(searchResult)
	
	// 4. save the posts to cache
	if data, err := json.Marshal(posts); err == nil {
		_ = backend.RedisBackend.Set(ctx, cacheKey, data, 10*time.Second)
	}

	// 5. construct response: researchResult -> []model.Post
	return posts, nil
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
	// GCS and ES are two independent services
	// So we need Saga 模式 —— 补偿事务（Compensating Transaction）
	// Normal Transaction: Save to GCS -> Save to ES -> Save to Redis
	// Abnormal Transaction: Save to ES failed -> Compensation: Delete the GCS file
	
	// 1. save the GCS and get URL
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.PostId)
	if err != nil {
		return err
	}

	post.Url = medialink
	post.Deleted = false
	post.DeletedAt = 0
	post.CleanupStatus = ""
	post.RetryCount = 0
	post.LastError = ""

	// 2. save the post to Elasticsearch
	if err := backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.PostId); err != nil {
		if deleteErr := backend.GCSBackend.DeleteFromGCS(post.Url); deleteErr != nil {
			fmt.Printf("CRITICAL: GCS orphan file, manual cleanup needed. PostId=%s, ESErr=%v, GCSErr=%v\n",
			post.PostId, err, deleteErr)
		}
		return fmt.Errorf("failed to save post to Elasticsearch: %v", err)
	}

	// 3. if both save to GCS and ES are successful, save to Redis
	ctx := context.Background()
	_ = backend.RedisBackend.Delete(ctx, utils.userFeedCacheKey(post.UserId))

	return nil
}

func DeletePost(postId, userId string) (bool, error) {
	if postId == "" || userId == "" {
		return false, nil
	}

	postQuery := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("post_id", postId)).
		MustNot(elastic.NewTermQuery("deleted", true))
	searchResult, err := backend.ESBackend.ReadFromES(postQuery, constants.POST_INDEX)
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
	
	err = backend.ESBackend.SaveToES(&post, constants.POST_INDEX, post.PostId)
	if err != nil {
		return false, err
	}

	err = backend.RedisBackend.Delete(ctx, utils.userFeedCacheKey(userId))
	if err != nil {
		return false, err
	}

	return true, nil
}

func LikePost(postId, userId string) (bool, error) {
	if postId == "" || userId == "" {
		return false, nil
	}
	ctx := context.Background()

	// 1. check whether the post exists
	postQuery := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("post_id", postId)).
		MustNot(elastic.NewTermQuery("deleted", true))
	searchResult, err := backend.ESBackend.ReadFromES(postQuery, constants.POST_INDEX)
	if err != nil {
		return false, err
	}
	posts := getPostFromSearchResult(searchResult)
	if len(posts) == 0 {
		return false, nil
	}
	post := posts[0]

	// 2. if exists, check whether the user already liked the post in Redis
	likeSetKey := fmt.Sprintf("like_set:%s", postId)
	alreadyLiked, err := backend.RedisBackend.SIsMember(ctx, likeSetKey, userId)
	if err == nil && alreadyLiked {
		return false, nil
	}

	// 3. if cache not hit, check whether the like exists in Elasticsearch
	likeId := postId + "_" + userId
	likeExistsQuery := elastic.NewTermQuery("post_like_id", likeId)
	likeExistsResult, err := backend.ESBackend.ReadFromES(likeExistsQuery, constants.LIKE_INDEX)
	if err != nil {
		return false, err
	}
	if likeExistsResult.TotalHits() > 0 {
		// 4. if like exists, update like in Redis, return true
		_ = backend.RedisBackend.SAdd(ctx, likeSetKey, userId)
		// Return true means the user create a new like, but the like already exists in Like Index in Elasticsearch, so return false
		return false, nil
	}

	// 4. if like not exists, create like and save to Like Index in Elasticsearch
	like := model.PostLike{
		PostLikeId: likeId,
		UserId: userId,
		PostId: postId,
		CreatedAt: time.Now().Unix(),
	}
	if err := backend.ESBackend.SaveToES(&like, constants.LIKE_INDEX, like.PostLikeId); err != nil {
		return false, err
	}

	// 5. update like count in Post Index in Elasticsearch
	post.LikeCount++
	if err := backend.ESBackend.IncrementFieldInES(constants.POST_INDEX, post.PostId, "like_count", 1); err != nil {
		return false, err
	}

	// 6. update like set in Redis
	_ = backend.RedisBackend.SAdd(ctx, likeSetKey, userId)
	return true, nil
}

func SharePost(postId, userId, platform string) (bool, error) {
	if postId == "" || userId == "" {
		return false, nil
	}

	postQuery := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("post_id", postId)).
		MustNot(elastic.NewTermQuery("deleted", true))
	searchResult, err := backend.ESBackend.ReadFromES(postQuery, constants.POST_INDEX)
	if err != nil {
		return false, err
	}
	posts := getPostFromSearchResult(searchResult)
	if len(posts) == 0 {
		return false, nil
	}
	post := posts[0]

	shareId := fmt.Sprintf("%s_%s_%s_%d", postId, userId, platform, time.Now().Unix())
	share := model.PostShare{
		PostShareId: shareId,
		UserId: userId,
		PostId: postId,
		CreatedAt: time.Now().Unix(),
		Platform: platform,
	}
	if err := backend.ESBackend.SaveToES(&share, constants.SHARE_INDEX, share.PostShareId); err != nil {
		return false, err
	}

	if err := backend.ESBackend.IncrementFieldInES(constants.POST_INDEX, post.PostId, "shared_count", 1); err != nil {
		return false, err
	}
	return true, nil
}

	