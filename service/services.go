package service

import (
	"github.com/iamrz1/ab-auth/config"
	"github.com/iamrz1/ab-auth/infra"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	"github.com/iamrz1/ab-auth/repo"
	rLog "github.com/iamrz1/rest-log"
)

// Config holds application configurations
type Config struct {
	CustomerService *customerService
	MerchantService *merchantService
}

// getServiceConfig returns service config
func getServiceConfig(cs *customerService, ms *merchantService) *Config {
	return &Config{CustomerService: cs, MerchantService: ms}
}

func SetupServiceConfig(cfg *config.AppConfig, db infra.DB, cache *infraCache.Redis, rLogger rLog.Logger) *Config {
	customerRepo := repo.NewCustomerRepo(db, cfg.CustomerTable, cache, rLogger)
	merchantRepo := repo.NewMerchantRepo(db, cfg.MerchantTable, cache, rLogger)
	addressRepo := repo.NewAddressRepo(db, cfg.AddressTable, "address_preset", rLogger)
	commonRepo := repo.NewCommonRepo(db, cache, rLogger)
	cs := NewCustomerService(cfg, commonRepo, customerRepo, addressRepo, rLogger)
	ms := NewMerchantService(cfg, commonRepo, merchantRepo, rLogger)

	return getServiceConfig(cs, ms)
}
