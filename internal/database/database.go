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
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    // Cloud SQL用の接続文字列フォーマットに変更
    var dsn string
    if strings.HasPrefix(dbHost, "/cloudsql/") {
        // Cloud SQL環境での接続文字列
        dsn = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbUser, dbPass, dbHost, dbName)
    } else {
        // ローカル環境での接続文字列
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            dbUser, dbPass, dbHost, dbPort, dbName)
    }

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    DB = db
    return nil
}

func GetDB() *gorm.DB {
    return DB
}