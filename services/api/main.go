package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MDmitryM/async-order-system/services/api/handler"
	"github.com/MDmitryM/async-order-system/services/api/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Println("api service")
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	env := os.Getenv("ENV")
	if env != "prod" {
		if err := godotenv.Load("./services/api/.env"); err != nil {
			logrus.Fatalf("error while reding .env, %s", err.Error())
		}
	}

	cfg := repository.PostresConfig{
		Host:     os.Getenv("API_DB_HOST"),
		Port:     os.Getenv("API_DB_PORT"),
		PG_User:  os.Getenv("POSTGRES_USER"),
		PG_Pwd:   os.Getenv("POSTGRES_PASSWORD"),
		PG_DB:    os.Getenv("POSTGRES_DB"),
		SSL_Mode: os.Getenv("API_DB_SSL_MODE"),
	}

	rootCtx := context.Background()

	pgPool, err := repository.NewPostgresDB(rootCtx, cfg)
	if err != nil {
		logrus.Fatalf("Error while creating connection pool: %v", err.Error())
	}
	defer pgPool.Close()

	app := fiber.New(fiber.Config{
		//Prefork: true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Order system v0.0",
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  10 * time.Second,
	})

	h := handler.NewHandler(&pgPool)
	h.InitRouts(app)

	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			logrus.Fatalf("error while server start, %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := app.Shutdown(); err != nil {
		logrus.Fatalf("error while server shutdown, %s", err.Error())
	}

	logrus.Println("server gracefully stopped!")
}
