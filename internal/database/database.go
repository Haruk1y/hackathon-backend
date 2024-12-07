package database

import (
    "fmt"
    "os"
	"strings"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    var dsn string
    if strings.HasPrefix(dbHost, "/cloudsql/") {
        // Cloud SQL Unix socket接続
        dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbUser, dbPass, dbHost, dbName)
    } else {
        // ローカル開発環境用のTCP接続
        dbPort := os.Getenv("DB_PORT")
        if dbPort == "" {
            dbPort = "3306" // デフォルトポート
        }
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbUser, dbPass, dbHost, dbPort, dbName)
    }

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("failed to initialize database, got error %v", err)
    }

    // コネクションプールの設定
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %v", err)
    }

    // コネクションプールの設定
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetMaxOpenConns(10)
    sqlDB.SetConnMaxLifetime(time.Hour)

    DB = db
    return nil
}

func GetDB() *gorm.DB {
    return DB
}