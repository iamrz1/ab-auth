package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConfig(t *testing.T) {
	err := LoadConfig()
	assert.NoError(t, err, "failed to load config")

	cfg := GetConfig()
	assert.NotEqualValues(t, cfg, nil, "failed to fetch config")

}
