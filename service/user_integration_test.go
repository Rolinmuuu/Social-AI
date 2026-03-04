//go:build integration
package service

import (
	"os"
	"testing"
	"socialai/backend"
	"socialai/model"
)

func TestMain(m *testing.M) {
	backend.InitElasticsearchBackend()
	os.Exit(m.Run())
}

func TestAddUser_Integration(t *testing.T) {
	testUsername := "testuser_integration"
	t.Cleanup(func() {
		err := backend.ESBackend.DeleteFromES(testUsername, constants.USER_INDEX)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
	success, err := AddUser(&model.User{
		Username: testUsername,
		Password: "testpassword",
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !success {
		t.Errorf("Expected true, got false")
	}
}