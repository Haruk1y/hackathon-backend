package handler

import (
    "gorm.io/gorm"
    "github.com/Haruk1y/hackathon-backend/internal/database"
)

var db *gorm.DB

// InitHandler initializes the handler package
func InitHandler() {
    db = database.GetDB()
}