package handler

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"regexp"

	"socialai/model"
	"socialai/service"
	
	"github.com/form3tech-oss/jwt-go"

)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signup request")
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		http.Error(w, "Request body is required", http.StatusBadRequest)
		return
	}

	// process request: JSON -> model.User 
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Failed to detect data from client", http.StatusBadRequest)
		fmt.Printf("Failed to detect data from client: %v\n", err)
		return
	}
	// validation
	if user.UserId == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9_]+$`).MatchString(user.UserId) == false {
		http.Error(w, "Invalid user_id or password", http.StatusBadRequest)
		fmt.Printf("Invalid user_id or password\n")
		return
	}

	// Call service layer: AddUser
	err := service.AddANewUser(&user)
	if err != nil {
		http.Error(w, "Failed to add user to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to add user to backend: %v\n", err)
		return
	}
	if !success {
		http.Error(w, "User already exists", http.StatusBadRequest)
		fmt.Printf("User already exists\n")
		return
	}

	// response
	fmt.Printf("User is added successfully! UserId: %s\n", user.UserId)
}

var mySigninKey = []byte("secret223344")
func signinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one signin request")
	w.Header().Set("Content-Type", "application/json")
	if r.Body == nil {
		http.Error(w, "Request body is required", http.StatusBadRequest)
		return
	}
	
	// process request: JSON -> model.User
	decoder := json.NewDecoder(r.Body)
	var user model.User
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, "Failed to detect data from client", http.StatusBadRequest)
		fmt.Printf("Failed to detect data from client: %v\n", err)
		return
	}

	// Call service layer: CheckUser
	err := service.CheckUserAlreadyExisted(user.UserId, user.Password)
	if err != nil {
		http.Error(w, "Failed to check user in backend", http.StatusInternalServerError)
		fmt.Printf("Failed to check user in backend: %v\n", err)
		return
	}
	if !success {
		http.Error(w, "Invalid user_id or password", http.StatusBadRequest)
		fmt.Printf("Invalid user_id or password\n")
		return
	}

	// response + generate token (JWT)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(mySigninKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		fmt.Printf("Failed to generate token: %v\n", err)
		return
	}
	w.Write([]byte(tokenString))
}