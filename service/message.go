package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"socialai/model"
	"socialai/backend"
	"socialai/constants"

	"github.com/google/uuid"
	elastic "github.com/olivere/elastic/v7"
)

func SendMessage(senderId, receiverId, content string) (string, error) {
	
	if senderId == receiverId {
		return "", fmt.Errorf("cannot send message to yourself")
	}

	userQuery := elastic.NewBoolQuery()
	.Filter(elastic.NewTermQuery("user_id", senderId))
	.Filter(elastic.NewTermQuery("user_id", receiverId))
	searchResult, err := backend.ESBackend.ReadFromES(userQuery, constants.USER_INDEX)
	if err != nil {
		return "", err
	}
	if searchResult.TotalHits() == 0 {
		return "", fmt.Errorf("user not found")
	}
	
	message := model.Message{
		MessageId: uuid.New().String(),
		SenderId: senderId,
		ReceiverId: receiverId,
		Content: content,
		CreatedAt: time.Now(),
	}
	_, err = backend.ESBackend.SaveToES(message, constants.MESSAGE_INDEX, messageId)
	if err != nil {
		return "", err
	}
	return messageId, nil
}

