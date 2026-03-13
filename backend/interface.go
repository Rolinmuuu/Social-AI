package backend

import (
	"io"
	"github.com/olivere/elastic/v7"
)

type ElasticsearchBackendInterface interface {
	ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error)
	SaveToES(i interface{}, index string, id string) error
	DeleteFromES(index string, id string) (bool, error)
	IncrementFieldInES(index string, id string, field string, value int) error
}

type GoogleCloudStorageBackendInterface interface {
	SaveToGCS(r io.Reader, objectName string) (string, error)
	DeleteFromGCS(objectName string) error
}

type RedisBackendInterface interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	SAdd(key string, members ...interface{}) error
	SIsMember(key string, member interface{}) (bool, error)
}