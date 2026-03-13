package utils

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"
)

const (
	// 4xx errors 客户端错误
	ErrCodeBadRequest = 40000 // 请求错误
	ErrCodeUnauthorized = 40100 // 未授权
	ErrCodeUserAlreadyExisted = 40001 // 用户已存在
	ErrCodeInvalidUser = 40002 // 用户/密码错误
	ErrCodeForbidden = 40300 // 禁止访问
	ErrCodeNotFound = 40400 // 未找到
	ErrCodePostNotFound = 40401 // 帖子未找到
	ErrCodeCommentNotFound = 40402 // 评论未找到
	ErrCodeLikeNotFound = 40403 // 点赞未找到
	ErrCodeShareNotFound = 40404 // 分享未找到
	ErrCodeFollowNotFound = 40405 // 关注未找到
	ErrCodeMessageNotFound = 40406 // 消息未找到
	ErrCodeUserNotFound = 40407 // 用户未找到
	ErrCodeAlreadyLiked = 40408 // 已经点赞
	ErrCodeAlreadyFollowed = 40409 // 已经关注
	ErrCodeMethodNotAllowed = 40500 // 方法不允许

	// 5xx errors 服务器错误
	ErrCodeInternalServerError = 50000 // 服务器错误
	ErrCodeESFailed = 50001 // Elasticsearch 操作失败
	ErrCodeRedisFailed = 50002 // Redis 操作失败
	ErrCodeGCSFailed = 50003 // GCS 操作失败
)

type APIResponse struct {
	RequestId string `json:"request_id"`
	Error *APIError `json:"error,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type APIError struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func WriteSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		RequestId: uuid.New().String(),
		Data: data,
	})
}

func WriteError(w http.ResponseWriter, httpStatus, errCode int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(APIResponse{
		RequestId: uuid.New().String(),
		Error: &APIError{
			Code: errCode,
			Message: errMessage,
		},
	})
}