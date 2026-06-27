package main

import (
	"database/sql"
	"log"
	"net/http"

	"go-first-api/internal/config"
	"go-first-api/internal/infrastructure/database"
	"go-first-api/internal/infrastructure/middleware"
	"go-first-api/internal/modules/auth"
	"go-first-api/internal/modules/health"
	"go-first-api/internal/modules/post"
	"go-first-api/internal/modules/user"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := sql.Open("pgx", cfg.DB.DSN())
	if err != nil {
		log.Fatal("failed to open database: ", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatal("failed to run migrations: ", err)
	}

	db, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	userRepo := user.NewRepository(db)
	userHandler := user.NewHandler(user.NewService(userRepo))

	postHandler := post.NewHandler(post.NewService(post.NewRepository(db)))

	authHandler := auth.NewHandler(auth.NewService(userRepo, cfg.JWT.Secret))

	healthHandler := health.NewHandler()

	jwtMiddleware := middleware.JWT(userRepo, cfg.JWT.Secret)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})

	healthHandler.RegisterRoutes(r)
	authHandler.RegisterRoutes(r, jwtMiddleware)
	userHandler.RegisterRoutes(r, jwtMiddleware)
	postHandler.RegisterRoutes(r, jwtMiddleware)

	log.Fatal(r.Run(cfg.Server.Addr))
}
