package main

import (
   "fmt"
   "log"
   "time"
   "net/http" 

   "socialai/backend"
   "socialai/service"
   "socialai/handler"   
)
func main() {
   fmt.Println("started-service")

   esBackend, err := backend.InitElasticsearchBackend()
   if err != nil {
      log.Fatal(err)
   }
   backend.ESBackend = esBackend
   gcsBackend, err := backend.InitGCSBackend()
   if err != nil {
      log.Fatal(err)
   }
   backend.GCSBackend = gcsBackend

   go func() {
      ticker := time.NewTicker(10 * time.Second)
      defer ticker.Stop()
      for range ticker.C {
         if success, err := service.CleanupDeletedPost(10); err != nil {
            log.Printf("Failed to cleanup deleted posts: %v", err)
         }
      }
   }()
   
   log.Fatal(http.ListenAndServe(":8080", handler.InitRouter()))
}