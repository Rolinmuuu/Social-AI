package handler

import (
	"net/http"
	
	"socialai/service"
	"encoding/json"
	"github.com/gorilla/mux"
)

type addCommentRequest struct {
	Content string `json:"content"`
	ParentCommentId string `json:"parent_comment_id"`
}

func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	if postId == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	userId, err := utils.GetUserIdFromJwtToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var addCommentRequest addCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&addCommentRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commentId, err := service.AddComment(postId, addCommentRequest.ParentCommentId, userId, addCommentRequest.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"comment_id": commentId})	
}