package service

import (
	"socialai/model"
	"socialai/backend"
	"socialai/constants"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func CheckUserAlreadyExisted(userId, password string) error {
	// Logic: search user_id+password in ES, if exists, return true, otherwise return false.
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("user_id", userId))
	query.Must(elastic.NewTermQuery("password", password))
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return fmt.Errorf("failed to read user from ES: %v", err)
	}
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Login as: %s\n", userId)
		return nil
	}
	return ErrUserNotFound
}

func AddANewUser(user *model.User) error {
	// 1. logic, check user id, then create
	query := elastic.NewTermQuery("user_id", user.UserId)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return fmt.Errorf("failed to read user from ES: %v", err)
	}
	if searchResult.TotalHits() > 0 {
		return ErrUserAlreadyExisted
	}
	
	// 2. ES save
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.UserId)
	if err != nil {
		return fmt.Errorf("failed to save user to ES: %v", err)
	}
	
	// 3. return
	fmt.Printf("User is added successfully! UserId: %s\n", user.UserId)
	return nil
}

