package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-user-api/internal/user"

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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// AutoMigrate — аналог synchronize: true в TypeORM
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("failed to migrate: ", err)
	}

	repo := user.NewRepository(db)
	service := user.NewService(repo)
	handler := user.NewHandler(service)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})
	handler.RegisterRoutes(r)

	log.Fatal(r.Run(":3000"))
}