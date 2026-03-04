package service

import (
	"testing"
	"socialai/model"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

type MockBackend struct {
	ReadFromESFunc func(query elastic.Query, index string) (*elastic.SearchResult, error)
	SaveToESFunc func(i interface{}, index string, id string) error
}

func (backend *MockBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	return backend.ReadFromESFunc(query, index)
}

func (backend *MockBackend) SaveToES(i interface{}, index string, id string) error {
	return backend.SaveToESFunc(i, index, id)
}

func TestUserService_UserAlreadyExists(t *testing.T) {
	success, err := AddUser(&model.User{
		Username: "testuser",
		Password: "testpassword",
	}, &MockBackend{
		ReadFromESFunc: func(query elastic.Query, index string) (*elastic.SearchResult, error) {
			return &elastic.SearchResult{
				TotalHits: func() int64 { return 1 },
			}, nil
		},
		SaveToESFunc: func(i interface{}, index string, id string) error {
			return nil
		}
	})
	assert.NoError(t, err)
	assert.True(t, success)
}