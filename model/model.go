package model

type Post struct {
   Id      string `json:"id"`
   User    string `json:"user"`
   Message string `json:"message"`
   Url     string `json:"url"`
   Type    string `json:"type"`
   Deleted bool   `json:"deleted"`
   DeletedAt int64 `json:"deleted_at"`
   CleanupStatus string `json:"cleanup_status"`
   RetryCount int `json:"retry_count"`
   LastError string `json:"last_error"`
}

type User struct {
   Username string `json:"username"`
   Password string `json:"password"`
   Age      int64  `json:"age"`
   Gender   string `json:"gender"`
}

