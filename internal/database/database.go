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
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbHost := os.Getenv("DB_HOST")
    
    var dsn string
    if os.Getenv("ENV") == "production" {
        // Cloud SQL用の接続文字列
        dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbUser, dbPass, dbHost, dbName)
    } else {
        // ローカル開発用の接続文字列
        dbPort := os.Getenv("DB_PORT")
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbUser, dbPass, dbHost, dbPort, dbName)
    }

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    if err := db.AutoMigrate(
        &model.User{},
        &model.Post{},
        &model.Like{},
    ); err != nil {
        return err
    }

    DB = db
    return nil
}

func GetDB() *gorm.DB {
    return DB
}