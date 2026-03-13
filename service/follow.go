package service

import (
	"fmt"
	"time"

	"socialai/backend"
	"socialai/constants"
	"socialai/model"

	elastic "github.com/olivere/elastic/v7"
	"github.com/google/uuid"
)

func AddFollow(followerId, followeeId string) (string, error) {
	if followerId == followeeId {
		return "", fmt.Errorf("cannot follow yourself")
	}

	followQuery := elastic.NewBoolQuery().
		.Filter(elastic.NewTermQuery("follower_id", followerId)).
		.Filter(elastic.NewTermQuery("followee_id", followeeId))
	searchResult, err := backend.ESBackend.ReadFromES(followQuery, constants.FOLLOW_INDEX)
	if err != nil {
		return "", err
	}
	if searchResult.TotalHits() > 0 {
		return "", fmt.Errorf("already following")
	}

	follow := model.Follow{
		FollowId: uuid.New().String(),
		FollowerId: followerId,
		FolloweeId: followeeId,
		CreatedAt: time.Now(),
	}

	_, err = backend.ESBackend.SaveToES(follow, constants.FOLLOW_INDEX, follow.FollowId)
	if err != nil {
		return "", err
	}

	return follow.FollowId, nil
}

func RemoveFollow(followerId, followeeId string) error {
	followQuery := elastic.NewBoolQuery().
		.Filter(elastic.NewTermQuery("follower_id", followerId)).
		.Filter(elastic.NewTermQuery("followee_id", followeeId))
	searchResult, err := backend.ESBackend.ReadFromES(followQuery, constants.FOLLOW_INDEX)
	if err != nil {
		return err
	}
	if searchResult.TotalHits() == 0 {
		return fmt.Errorf("not following")
	}

	var followId string
	for _, hit := range searchResult.Hits.Hits {
		followId = hit.Id
		break
	}

	_, err = backend.ESBackend.DeleteFromES(constants.FOLLOW_INDEX, followId)
	if err != nil {
		return err
	}

	return nil
}

func GetFollowerIds(followerId string) ([]string, error) {
	followQuery := elastic.NewTermQuery("follower_id", followerId)
	searchResult, err := backend.ESBackend.ReadFromES(followQuery, constants.FOLLOW_INDEX)
	if err != nil {
		return nil, err
	}

	var followerIds []string
	for _, hit := range searchResult.Hits.Hits {
		var follow model.Follow
		if err := json.Unmarshal(hit.Source, &follow); err == nil {
			followerIds = append(followerIds, follow.FollowerId)
		}
	}

	return followerIds, nil
}