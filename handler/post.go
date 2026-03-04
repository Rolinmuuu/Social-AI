package handler

import (
   "encoding/json"
   "fmt"
   "net/http"
   "path/filepath"
   "socialai/model"
   "socialai/service"

   jwt "github.com/form3tech-oss/jwt-go"
   "github.com/gorilla/mux"
   "github.com/pborman/uuid"
)

var (
    mediaTypes = map[string]string{
        ".jpg": "image",
        ".jpeg": "image",
        ".gif": "image",
        ".png": "image",
        ".mp4": "video",
        ".avi": "video",
        ".mov": "video",
        ".flv": "video",
        ".wmv": "video",
    }
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
   // Parse from body of request to get a json object.
   fmt.Println("Received one upload request")

    // 1. Process POST request: 
    // multipart request (JSON String) -> model.Post, Image/Video (Post struct).
    // token -> username
    token := r.Context().Value("user")
    claims := token.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

    p := model.Post{
        Id: uuid.New(),
        User: username.(string),
        Message: r.FormValue("message"),
    }
    file, header, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Failed to read media file", http.StatusBadRequest)
        fmt.Printf("Failed to read media file: %v\n", err)
        return
    }
    suffix := filepath.Ext(header.Filename)
    if mediaType, ok := mediaTypes[suffix]; ok {
        p.Type = mediaType
    } else {
        p.Type = "unknown"
    }


   
    // 2. Call Service layer (Business logic) to save the post to database.
    err = service.SavePost(&p, file)
    if err != nil {
        http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
        fmt.Printf("Failed to save post to backend: %v\n", err)
        return
    }



	// 3. Construct response and send back to client.
   fmt.Printf("Post is saved successfully!\n")
   fmt.Fprintf(w, "Post is saved successfully!")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
   fmt.Println("Received one search request")
   // 1. process the request: URL params => user + keywords
   user := r.URL.Query().Get("user")
   keywords := r.URL.Query().Get("keywords")

   // 2. call service layer to search posts in database. (business logic)
   var posts []model.Post
   var err error
   if user != "" {
       posts, err = service.SearchPostByUser(user)
   } else {
       posts, err = service.SearchPostByKeywords(keywords)
   }
   if err != nil {
       http.Error(w, "Failed to read from backend", http.StatusInternalServerError)
       fmt.Printf("Failed to read from backend: %v\n", err)
       return
   }
   
   // 3. construct response and send back to client.
   // []model.Post -> JSON String
   js, err := json.Marshal(posts)
   if err != nil {
       http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
       fmt.Printf("Failed to parse posts into JSON format: %v\n", err)
       return
   }
   w.Header().Set("Content-Type", "application/json")
   w.Write(js)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one delete request")
	postID := mux.Vars(r)["id"]
	if postID == "" {
		http.Error(w, "Missing post id", http.StatusBadRequest)
		return
	}

	deleted, err := service.DeletePost(postID)
	if err != nil {
		http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to delete post from backend: %v\n", err)
		return
	}
	if !deleted {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Post is marked as deleted, cleanup in progress"))
}