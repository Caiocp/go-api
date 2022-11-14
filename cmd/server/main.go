package main

import (
	"net/http"

	"github.com/caiocp/go-api/configs"
	"github.com/caiocp/go-api/internal/entities"
	"github.com/caiocp/go-api/internal/infra/database"
	"github.com/caiocp/go-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entities.Product{}, &entities.User{})

	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)

	productHandler := handlers.NewProductHandler(productDB)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", configs.JwtExpiresIn))
	// r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	http.ListenAndServe(":8080", r)
}

// func LogRequest(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Printf("%s %s", r.Method, r.URL.Path)
// 		next.ServeHTTP(w, r)
// 	})
// }
