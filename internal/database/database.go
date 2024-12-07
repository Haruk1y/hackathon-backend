package database

import (
    "fmt"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/Haruk1y/hackathon-backend/internal/model"
)

var DB *gorm.DB

func InitDB() error {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPass, dbHost, dbPort, dbName)

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    // Auto-migrate all models
    if err := db.AutoMigrate(
        &model.User{},
        &model.Post{},
        &model.Like{},  // Likeモデルを追加
    ); err != nil {
        return err
    }

    DB = db
    return nil
}

func GetDB() *gorm.DB {
    return DB
}