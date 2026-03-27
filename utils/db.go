package utils

import (
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConnector() *gorm.DB {
	dsn := strings.TrimSpace(os.Getenv("VMUSIC_DB_DSN"))
	if dsn == "" {
		dsn = "host=192.168.81.101 user=vmusic password=123456 dbname=vmusic_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("建立连接失败：%v", err)
	}
	return db
}
