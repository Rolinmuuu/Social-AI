package backend

import (
	"os"
	"strconv"
	"testing"
	"socialai/constants"
)

func TestMain(m *testing.M) {
	_, _ = InitElasticsearchBackend()
	os.Exit(m.Run())
}

func TestESBackend_SaveToES(t *testing.T) {
	testID := "integration-test-id-" + strconv.Itoa(os.Getpid())

	t.Cleanup(func() {
		_, _ = ESBackend.DeleteFromES(constants.USER_INDEX, testID)
	})

	err := ESBackend.SaveToES(map[string]interface{}{
		"user_id": testID,
	}, constants.USER_INDEX, testID)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}