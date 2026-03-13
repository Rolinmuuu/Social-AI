package backend

import (
   "context"
   "fmt"
   "socialai/constants"

   "github.com/olivere/elastic/v7"

)

var (
	ESBackend ElasticsearchBackendInterface
)

type ElasticsearchBackend struct {
	client *elastic.Client
}
func InitElasticsearchBackend() (ElasticsearchBackendInterface, error) {
   client, err := elastic.NewClient(
       elastic.SetURL(constants.ES_URL),
       elastic.SetBasicAuth(constants.ES_USERNAME, constants.ES_PASSWORD))
   if err != nil {
       return nil, err
   }

   exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
   if err != nil {
       panic(err)
   }

   if !exists {
       mapping := `{
           "mappings": {
               "properties": {
                   "post_id":  { "type": "keyword" },
                   "user_id":  { "type": "keyword" },
                   "user":     { "type": "keyword" },
                   "message":  { "type": "text" },
                   "url":      { "type": "keyword", "index": false },
                   "type":     { "type": "keyword", "index": false },
                   "deleted":  { "type": "boolean" },
                   "deleted_at": { "type": "long" },
                   "cleanup_status": { "type": "keyword" },
                   "retry_count": { "type": "integer" },
                   "last_error": { "type": "text" },
                   "like_count": { "type": "integer" },
                   "shared_count": { "type": "integer" }
               }
           }
       }`
       _, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
       if err != nil {
           panic(err)
       }
   }

   exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
   if err != nil {
       panic(err)
   }

   if !exists {
       mapping := `{
                       "mappings": {
                               "properties": {
                                        "user_id": {"type": "keyword"},
                                       "username": {"type": "keyword"},
                                       "password": {"type": "keyword"},
                                       "age":      {"type": "long", "index": false},
                                       "gender":   {"type": "keyword", "index": false}

                               }
                       }
               }`
       _, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
       if err != nil {
           panic(err)
       }
   }
   fmt.Println("Indexes are created.")


   exists, err = client.IndexExists(constants.FOLLOW_INDEX).Do(context.Background())
   if err != nil {
       panic(err)
   }
   if !exists {
       mapping := `{
           "mappings": {
               "properties": {
                   "follow_id": { "type": "keyword" },
                   "follower_id": { "type": "keyword" },
                   "followee_id": { "type": "keyword" },
                   "created_at": { "type": "long" }
               }
           }
       }`
       _, err = client.CreateIndex(constants.FOLLOW_INDEX).Body(mapping).Do(context.Background())
       if err != nil {
           panic(err)
       }
   }
   fmt.Println("Indexes are created.")

   ESBackend = &ElasticsearchBackend{client: client}
   return ESBackend, nil
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
   searchResult, err := backend.client.Search().
        Index(index).
        Query(query).
        Pretty(true).
        Do(context.Background())
   if err != nil {
       return nil, err
   }
   return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error {
    _, err := backend.client.Index().
        Index(index).
        Id(id).
        BodyJson(i).
        Do(context.Background())
    return err
}

func (backend *ElasticsearchBackend) DeleteFromES(index string, id string) (bool, error) {
	resp, err := backend.client.Delete().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return false, err
	}
	return resp.Result == "deleted", nil
}

func (backend *ElasticsearchBackend) IncrementFieldInES(index, id, field string, value int) error {
    scriptSource := "ctx._source[params.field] = (ctx._source[params.field] == null ? 0 : ctx._source[params.field]) + params.value"
    script := elastic.NewScript(scriptSource).Params(map[string]interface{}{"field": field, "value": value})
    return, err := backend.client.Update().
        Index(index).
        Id(id).
        Script(script).
        RetryOnConflict(3).
        Do(context.Background())
    if err != nil {
        return err
    }
    if result.Result != "updated" && result.Result != "noop" {
        return fmt.Errorf("increment failed: %s", result.Result)
    }
    return nil
}

func (backend *ElasticsearchBackend) GetByIdFromES(index string, id string) (interface{}, error) {
    get, err := backend.client.Get().
        Index(index).
        Id(id).
        Do(context.Background())
    if err != nil {
        return nil, err
    }
    return get.Source, nil
}

