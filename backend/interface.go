package backend

import (
	"io"
	"github.com/olivere/elastic/v7"
)

type ElasticsearchBackendInterface interface {
	ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error)
	SaveToES(i interface{}, index string, id string) error
	DeleteFromES(index string, id string) (bool, error)
}

type GoogleCloudStorageBackendInterface interface {
	SaveToGCS(r io.Reader, objectName string) (string, error)
	DeleteFromGCS(objectName string) error
}