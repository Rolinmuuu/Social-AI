package service

import (
	"socialai/model"
	"socialai/backend"
	"socialai/constants"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func CheckUser(username, password string) (bool, error) {
	// Logic: search user+password in ES, if exists, return true, otherwise return false.
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("username", username))
	query.Must(elastic.NewTermQuery("password", password))
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}
	if searchResult.TotalHits() > 0 {
		fmt.Printf("Login as: %s\n", username)
		return true, nil
	}
	return false, nil
}

func AddUser(user *model.User) (bool, error) {
	// 1. logic, check username, then create
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}
	if searchResult.TotalHits() > 0 {
		return false, nil
	}
	
	// 2. ES save
	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	
	// 3. return
	fmt.Printf("User is added successfully! Username: %s\n", user.Username)
	return true, err
}

