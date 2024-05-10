package database

import (
	"database/sql"
	"os"
	"fmt"
	"sync"

	"example.com/server/types"
)

type Database struct {
	*sql.DB
}

var (
	dbInstance *Database
	once sync.Once
)

type Config struct {
	User string
	Password string
	Host string
	Port string
	Name string
}

type DBLogger struct{}

func (l *DBLogger) Info(message string) {
	fmt.Println("INFO:", message)
}

func (l *DBLogger) Error(message string) {
	fmt.Println("ERROR:", message)
}

var dblogger types.Logger = &DBLogger{}

func getConfig() Config {
	return Config {
		User: getEnv("DBUSER", "defaultUser"),
		Password: getEnv("DBPASS", "defaultPassword"),
		Host: getEnv("DBHOST", "localhost"),
		Port: getEnv("DBPORT", "3306"),
		Name: getEnv("DBNAME", "defaultDB"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func connectDB() (*Database, error) {
	config := getConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.Name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to conned to database: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dblogger.Info("Connected to database")
	return &Database{DB: db}, nil
}

func GetDB() (*Database, error) {
	var err error
	once.Do(func() {
		dbInstance, err = connectDB()
	})
	return dbInstance, err
}