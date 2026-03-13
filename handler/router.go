package handler

import (
   "net/http" 

   "socialai/middleware"

   jwtMiddleware "github.com/auth0/go-jwt-middleware"
   jwt "github.com/form3tech-oss/jwt-go"

   "github.com/gorilla/handlers"
   "github.com/gorilla/mux"   
   "github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitRouter() http.Handler {
   jwtAuth := jwtMiddleware.New(jwtMiddleware.Options{
       ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
           return []byte(mySigninKey), nil
       },
       SigningMethod: jwt.SigningMethodHS256,
   })
      

   router := mux.NewRouter()

   // metrics routes
   router.Handle("/metrics", promhttp.Handler()).Methods("GET")
   // give the metrics middleware to all routes
   router.Use(middleware.MetricsMiddleware)
   router.Use(middleware.LoggingMiddleware)

   // auth routes
   router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
   router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")
  
   // post routes
   router.Handle("/upload", jwtAuth.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
   router.Handle("/search", jwtAuth.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
   router.Handle("/post/{id}", jwtAuth.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")
   router.Handle("/post/{id}/like", jwtAuth.Handler(http.HandlerFunc(likeHandler))).Methods("POST")
   router.Handle("/post/{id}/share", jwtAuth.Handler(http.HandlerFunc(shareHandler))).Methods("POST")

   // comment routes
   router.Handle("/post/{id}/comment", jwtAuth.Handler(http.HandlerFunc(addCommentHandler))).Methods("POST")

   // follow routes
   router.Handle("/follow", jwtAuth.Handler(http.HandlerFunc(addFollowHandler))).Methods("POST")
   router.Handle("/follow", jwtAuth.Handler(http.HandlerFunc(removeFollowHandler))).Methods("DELETE")
   router.Handle("/follow/followers", jwtAuth.Handler(http.HandlerFunc(getFollowerIdsHandler))).Methods("GET")

   // message routes
   router.Handle("/message", jwtAuth.Handler(http.HandlerFunc(sendMessageHandler))).Methods("POST")
   router.Handle("/message", jwtAuth.Handler(http.HandlerFunc(getMessageHandler))).Methods("GET")
   
   origins := handlers.AllowedOrigins([]string{"*"})
   methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
   headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

   return handlers.CORS(origins, methods, headers)(router)
}
