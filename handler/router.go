package handler

import (
   "net/http" 


   jwtMiddleware "github.com/auth0/go-jwt-middleware"

   jwt "github.com/form3tech-oss/jwt-go"

   "github.com/gorilla/handlers"
   "github.com/gorilla/mux"   
)

func InitRouter() http.Handler {
   middleware := jwtMiddleware.New(jwtMiddleware.Options{
       ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
           return []byte(mySigninKey), nil
       },
       SigningMethod: jwt.SigningMethodHS256,
   })
      

   router := mux.NewRouter()

   router.Handle("/upload", middleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
   router.Handle("/search", middleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
   router.Handle("/post/{id}", middleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")
   router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
   router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")

   origins := handlers.AllowedOrigins([]string{"*"})
   methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
   headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

   return handlers.CORS(origins, methods, headers)(router)
}
