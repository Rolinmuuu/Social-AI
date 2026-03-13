package utils

import (
	"context"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)

func GetUserIdFromJwtToken(r *http.Reques) (string, error) {
	token := r.Context().Value("user")
    jwtToken, ok := token.(*jwt.Token)
    if !ok || jwtToken == nil {
        return "", errors.New("Unauthorized: Missing or invalid token")
    }
    claims, ok := jwtToken.Claims.(jwt.MapClaims)
    if !ok {
        return "", errors.New("Unauthorized: Missing or invalid claims")
    }
    userId, ok := claims["user_id"].(string)
    if !ok || userId == "" {
        return "", errors.New("Unauthorized: Missing or invalid user id in token")
    }
    return userId, nil
}

func userFeedCacheKey(userId string) string {
	return fmt.Sprintf("user_feed:%s", userId)
}