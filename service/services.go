package service

import (
	"github.com/iamrz1/ab-auth/config"
	"github.com/iamrz1/ab-auth/infra"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/repo"
)

// Config holds application configurations
type Config struct {
	CustomerService *customerService
}

// getServiceConfig returns service config
func getServiceConfig(cs *customerService) *Config {
	return &Config{CustomerService: cs}
}

func SetupServiceConfig(cfg *config.AppConfig, db infra.DB, cache *infraCache.Redis, logStruct logger.StructLogger) *Config {
	customerRepo := repo.NewCustomerRepo(db, cfg.CustomerTable, cache, logStruct)
	addressRepo := repo.NewAddressRepo(db, cfg.AddressTable, "address_preset", logStruct)
	commonRepo := repo.NewCommonRepo(db, cache, logStruct)
	cs := NewCustomerService(cfg, commonRepo, customerRepo, addressRepo, logStruct)

	return getServiceConfig(cs)
}
