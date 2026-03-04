package backend

import (
	"os"
	"strconv"
	"time"
	"testing"
	"socialai/constants"
	"github.com/olivere/elastic/v7"
)

func TestMain(m *testing.M) {
	InitElasticsearchBackend()
	os.Exit(m.Run())
}

func TestESBackend_SaveToES(t *testing.T) {
	testID := "integration-test-id-" + strconv.Itoa(time.Now().Nanosecond())

	t.Cleanup(func() {
		ESBackend.DeleteFromES(elastic.NewTermQuery("id", testID), constants.USER_INDEX)
	})

	