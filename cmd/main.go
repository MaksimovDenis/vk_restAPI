package main

import (
	"fmt"
	"io"
	"os"
	filmoteke "vk_restAPI"
	logger "vk_restAPI/logs"
	"vk_restAPI/package/handler"
	"vk_restAPI/package/repository"
	"vk_restAPI/package/service"

	"gopkg.in/yaml.v3"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title HOCHU V VK
// @verstion 1.0
// @description API Server for Filmoteka Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @contact.name Denis Maksimov
// @contact.email maksimovis74@gmail.com

type Config struct {
	Port string `yaml:"port"`
	DB   struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBname   string `yaml:"dbname"`
		SSLmode  string `yaml:"sslmode"`
	}
}

func initConfig() (*Config, error) {
	var config Config

	file, err := os.Open("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %v", err)
	}

	return &config, nil
}

func main() {

	//Setting JSON format for our logs
	logrus.SetFormatter(new(logrus.JSONFormatter))

	//Loading .env
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	config, err := initConfig()
	if err != nil {
		logrus.Fatal("error initializing config:", err)
	}

	//Initializing our DB
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   config.DB.DBname,
		SSLMode:  config.DB.SSLmode,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	//Creating our dependencies
	repositories := repository.NewRepository(db)
	services := service.NewService(repositories)
	handlers := handler.NewHandler(services)

	logger.Log.Info("Server started")

	//Running server
	srv := new(filmoteke.Server)
	if err := srv.Run(config.Port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}

}
