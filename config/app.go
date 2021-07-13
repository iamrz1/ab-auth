package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

// AppConfig holds application configurations
type AppConfig struct {
	Environment     string
	Host            string
	Port            int
	OtpTtlMinutes int
	GracefulTimeout int
	DSN             string
	Database        string
	CustomerTable   string
	CacheURL string
}

var myConfig *AppConfig

func init() {
	godotenv.Load("../.env")
}

// LoadConfig loads config from path
func LoadConfig() error {
	port, err := strconv.Atoi(os.Getenv("REST_PORT"))
	if err != nil {
		log.Println("rest port not found or invalid")
		port = 8080
	}

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("missing env DB_URL")
	}

	dbname := os.Getenv("DB_DATABASE_NAME")
	if dbname == "" {
		log.Fatal("missing env DB_DATABASE_NAME")
	}

	ct := os.Getenv("DB_CUSTOMER_COLLECTION_NAME")
	if ct == "" {
		log.Fatal("missing env DB_CUSTOMER_COLLECTION_NAME")
	}

	cacheURL := os.Getenv("REDIS_URL")
	if cacheURL == "" {
		log.Fatal("missing env REDIS_URL")
	}

	otpttl, err := strconv.Atoi(os.Getenv("OTP_TTL_MINUTES"))
	if err != nil {
		log.Fatal("missing env OTP_TTL_MINUTES")
	}

	myConfig = &AppConfig{
		Environment:     os.Getenv("ENV"),
		Host:            os.Getenv("REST_HOST"),
		Port:            port,
		OtpTtlMinutes: otpttl,
		GracefulTimeout: 30,
		DSN:             dsn,
		Database:        dbname,
		CustomerTable:   ct,
		CacheURL: cacheURL,
	}

	return nil
}

// GetConfig returns application config
func GetConfig() *AppConfig {
	return myConfig
}
