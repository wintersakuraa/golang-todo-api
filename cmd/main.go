package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wintersakura/todo-api"
	"github.com/wintersakura/todo-api/pkg/handler"
	"github.com/wintersakura/todo-api/pkg/repository"
	"github.com/wintersakura/todo-api/pkg/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatal(err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatal(err)
	}

	db, err := repository.NewMysqlDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		Query:    viper.GetString("db.query"),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	s := new(todo.Server)
	go func() {
		if err := s.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatal(err)
		}
	}()

	logrus.Print("Todo API Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Todo API Shutting Down")

	if err := s.ShutDown(context.Background()); err != nil {
		logrus.Error(err)
	}
	if err := db.Close(); err != nil {
		logrus.Error(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
