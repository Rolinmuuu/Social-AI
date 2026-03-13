package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/form3tech-oss/jwt-go"
)

var routeManifest = []struct {
	method string
	path string
	protected bool
}{
	{"POST", "/signup", false},
	{"POST", "/signin", false},
	{"POST", "/upload", true},
	{"GET", "/search", true},
	{"DELETE", "/post/test-post-id", true},
}

func TestRouter_WrongMethod(t *testing.T) {
	router := InitRouter()
	allMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "CONNECT", "TRACE"}
	for _, route := range routeManifest {
		for _, wrongMethod := range allMethods {
			if wrongMethod == route.method {
				continue
			}
			req, err := http.NewRequest(wrongMethod, route.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			if rr.Code != http.StatusMethodNotAllowed {
				t.Errorf("[%s %s] Expected status code 405, got %d", wrongMethod, route.path, rr.Code)
			}
		}
	}
}

func TestRouter_AuthMiddleware(t *testing.T) {
	router := InitRouter()
	for _, route := range routeManifest {
		req, err := http.NewRequest(route.method, route.path, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if route.protected && rr.Code != http.StatusUnauthorized {
			t.Errorf("[%s %s] Expected status code 401, got %d", route.method, route.path, rr.Code)
		} 
		if !route.protected && rr.Code == http.StatusUnauthorized {
			t.Errorf("[%s %s] Public route should not return 401, but got %d", route.method, route.path, rr.Code)
		}
	}
}

func TestRouter_ValidToken_PassesMiddleware(t *testing.T) {
	router := InitRouter()

	validToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "test_user",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte(mySigninKey))
	if err != nil {
		t.Fatalf("Failed to create valid token: %v", err)
	}

	for _, route := range routeManifest {
		if route.protected {
			if route.path == "/search" || route.path == "/post/test-post-id" {
				// These handlers require initialized backends; middleware behavior is covered by other routes.
				continue
			}
			req, err := http.NewRequest(route.method, route.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Authorization", "Bearer "+validToken)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code == http.StatusUnauthorized {
				t.Errorf("[%s %s] Effective Token should not return 401, but got %d", route.method, route.path, rr.Code)
			}
		}
	}
}

