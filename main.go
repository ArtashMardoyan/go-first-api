package main

import (
	"log"
	"net/http"
	"os"

	"go-first-api/internal/auth"
	"go-first-api/internal/database"
	"go-first-api/internal/middleware"
	"go-first-api/internal/post"
	"go-first-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	if err := db.AutoMigrate(&user.User{}, &post.Post{}); err != nil {
		log.Fatal("failed to migrate: ", err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService)

	authService := auth.NewService(userRepo)
	authHandler := auth.NewHandler(authService)

	jwtMiddleware := middleware.JWT(userRepo)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})

	authHandler.RegisterRoutes(r, jwtMiddleware)
	userHandler.RegisterRoutes(r, jwtMiddleware)
	postHandler.RegisterRoutes(r, jwtMiddleware)

	log.Fatal(r.Run(":3000"))
}