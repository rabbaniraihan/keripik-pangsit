package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"keripik-pangsit/helper"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var DB *sqlx.DB
var Redis *redis.Client

func InitPostgres() {
	host := helper.GetEnv("DB_HOST", "localhost")
	port := helper.GetEnv("DB_PORT", "5432")
	user := helper.GetEnv("DB_USER", "postgres")
	password := helper.GetEnv("DB_PASSWORD", "postgres")
	dbname := helper.GetEnv("DB_NAME", "keripik_pangsit")
	sslmode := helper.GetEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	DB = db
	log.Println("PostgreSQL connected successfully")
}

func InitRedis() {
	addr := helper.GetEnv("REDIS_ADDR", "localhost:6379")
	password := helper.GetEnv("REDIS_PASSWORD", "")
	dbIndex := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       dbIndex,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	Redis = rdb
	log.Println("Redis connected successfully")
}
