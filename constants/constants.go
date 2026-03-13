package constants 

import "os"

const (
   USER_INDEX = "user"
   POST_INDEX = "post"
   LIKE_INDEX = "like"
   SHARE_INDEX = "share"
   COMMENT_INDEX = "comment"
   FOLLOW_INDEX = "follow"
   MESSAGE_INDEX = "message"

   REDIS_ADDRESS = "redis:6379"
   REDIS_PASSWORD = ""
   REDIS_DB = 0

   ES_URL = "http://10.128.0.2:9200"
   ES_USERNAME = "rolinmu"
   ES_PASSWORD = os.Getenv("ES_PASSWORD")
   GCS_BUCKET = "socialai_laioffer_202512"

   LOGSTASH_ADDRESS = "logstash:5000"
)
