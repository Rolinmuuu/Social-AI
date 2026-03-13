package handler

import (
	"net/http"
	"socialai/service"
)

func AddFollowHandler(w http.ResponseWriter, r *http.Request) {
	followerId, err := utils.GetUserIdFromJwtToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	followeeId := mux.Vars(r)["followee_id"]
	if followeeId == "" {
		http.Error(w, "Followee ID is required", http.StatusBadRequest)
		return
	}
	followId, err := service.AddFollow(followerId, followeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Followed successfully"})
}

func RemoveFollowHandler(w http.ResponseWriter, r *http.Request) {
	followerId, err := utils.GetUserIdFromJwtToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	followeeId := mux.Vars(r)["followee_id"]
	if followeeId == "" {
		http.Error(w, "Followee ID is required", http.StatusBadRequest)
		return
	}
	err = service.RemoveFollow(followerId, followeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Unfollowed successfully"})		
}

func GetFollowerIdsHandler(w http.ResponseWriter, r *http.Request) {
	followerId, err := utils.GetUserIdFromJwtToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	followerIds, err := service.GetFollowerIds(followerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string][]string{"follower_ids": followerIds})
}
