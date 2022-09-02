package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	todo "github.com/RomanKlyachev/go-rest-api"
	"github.com/RomanKlyachev/go-rest-api/pkg/handler"
	"github.com/RomanKlyachev/go-rest-api/pkg/repository"
	"github.com/RomanKlyachev/go-rest-api/pkg/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Initialization configs error: %s", err.Error())
	}

	if err := godotenv.Load("../.env"); err != nil {
		logrus.Fatalf("Loading ENV variable error: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     viper.GetString("DBHost"),
		Port:     viper.GetString("DBPort"),
		Username: viper.GetString("DBUsername"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("DBName"),
		SSLMode:  viper.GetString("DBSSLMode"),
	})
	if err != nil {
		logrus.Fatalf("Connection DB error: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Http server running error: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
