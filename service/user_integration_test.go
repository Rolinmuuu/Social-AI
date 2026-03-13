//go:build integration
package service

import (
	"os"
	"testing"
	"socialai/backend"
	"socialai/constants"
	"socialai/model"
)

func TestMain(m *testing.M) {
	_, _ = backend.InitElasticsearchBackend()
	os.Exit(m.Run())
}

func TestAddUser_Integration(t *testing.T) {
	testUserId := "testuser_integration"
	t.Cleanup(func() {
		_, err := backend.ESBackend.DeleteFromES(constants.USER_INDEX, testUserId)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	success, err := AddUser(&model.User{
		UserId:   testUserId,
		Username: testUserId,
		Password: "testpassword",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !success {
		t.Errorf("Expected true, got false")
	}
}