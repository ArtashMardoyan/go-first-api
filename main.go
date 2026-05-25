package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-first-api/internal/auth"
	"go-first-api/internal/post"
	"go-first-api/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	if err := db.AutoMigrate(&user.User{}, &post.Post{}); err != nil {
		log.Fatal("failed to migrate: ", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	userService := user.NewService(user.NewRepository(db))
	userHandler := user.NewHandler(userService)
	postHandler := post.NewHandler(post.NewService(post.NewRepository(db)))
	authHandler := auth.NewHandler(userService)

	jwtMiddleware := auth.JWTMiddleware(userService)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})
	userHandler.RegisterRoutes(r, jwtMiddleware)
	postHandler.RegisterRoutes(r, jwtMiddleware)
	authHandler.RegisterRoutes(r)

	log.Fatal(r.Run(":3000"))
}