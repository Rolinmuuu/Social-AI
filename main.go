package main

import (
   "fmt"
   "log"
   "time"
   "net/http" 

   "go.uber.org/zap"
   "socialai/backend"
   "socialai/service"
   "socialai/handler"   
   "socialai/constants"
   "socialai/logger"
)

func main() {
   logger.InitLogger(constants.LOGSTASH_ADDRESS)
   defer logger.Logger.Sync()

   fmt.Println("started-service")

   // Initialize Elasticsearch, GCS and Redis backends
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

   redisBackend, err := backend.InitRedisBackend()
   if err != nil {
      log.Fatal(err)
   }
   backend.RedisBackend = redisBackend

   // Start the cleanup process every 10 seconds
   go func() {
      ticker := time.NewTicker(10 * time.Second)
      defer ticker.Stop()
      for range ticker.C {
         if _, err := service.CleanupDeletedPost(10); err != nil {
            log.Printf("Failed to cleanup deleted posts: %v", err)
         }
      }
   }()
   
   if err := http.ListenAndServe(":8080", handler.InitRouter()); err != nil {
      logger.Logger.Error("Failed to start server", zap.Error(err))
   }
}