package handler

import (

	"testing"
)

func TestPostUploadHandler(t *testing.T) {
	if mediaTypes[".jpg"]!= "image" {
		t.Errorf("Expected image, got %s", mediaTypes[".jpg"])
	}
	if mediaTypes[".png"]!= "image" {
		t.Errorf("Expected image, got %s", mediaTypes[".png"])
	}
	if mediaTypes[".gif"]!= "image" {
		t.Errorf("Expected image, got %s", mediaTypes[".gif"])
	}
	if mediaTypes[".mp4"]!= "video" {
		t.Errorf("Expected video, got %s", mediaTypes[".mp4"])
	}
	if _, ok := mediaTypes[".unknown"]; ok {
		t.Errorf("Expected null, got %s", mediaTypes[".unknown"])
	} 
}