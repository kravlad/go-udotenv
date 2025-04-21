package udotenv

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, []string{"envs", "e"}, config.EnvFlags)
	assert.Equal(t, []string{"env-overload", "eo", "o"}, config.OverloadFlags)
	assert.Equal(t, defaultEnvPath, config.DefaultEnvPath)
	assert.False(t, config.OverloadByDefault)
}

func TestNew_DefaultConfig(t *testing.T) {
	udotEnv := New(false)

	assert.NotNil(t, udotEnv)
	assert.NotNil(t, udotEnv.Config)
	assert.Equal(t, defaultEnvPath, udotEnv.Config.DefaultEnvPath)
	assert.Empty(t, udotEnv.EnvParam)
	assert.False(t, udotEnv.OverloadParam)
}

func TestNew_CustomConfig(t *testing.T) {
	customConfig := &Config{
		EnvFlags:          []string{"custom-env"},
		OverloadFlags:     []string{"custom-overload"},
		DefaultEnvPath:    "custom.env",
		OverloadByDefault: true,
	}

	udotEnv := New(false, customConfig)

	assert.NotNil(t, udotEnv)
	assert.Equal(t, customConfig, udotEnv.Config)
	assert.Equal(t, "custom.env", udotEnv.Config.DefaultEnvPath)
	assert.True(t, udotEnv.Config.OverloadByDefault)
}

func TestNew_MultipleConfigsPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when passing multiple configs")
		}
	}()

	New(false, GetDefaultConfig(), GetDefaultConfig())
}

func TestLoad_NoEnvParam(t *testing.T) {
	udotEnv := &udotEnvType{}

	assert.NotPanics(t, func() {
		udotEnv.Load()
	})
}

func TestLoad_WithEnvParam(t *testing.T) {
	_ = godotenv.Write(map[string]string{"TEST_KEY": "TEST_VALUE"}, ".test.env")
	defer os.Remove(".test.env")

	udotEnv := &udotEnvType{
		EnvParam:      stringSlice{".test.env"},
		OverloadParam: false,
	}

	assert.NotPanics(t, func() {
		udotEnv.Load()
	})

	assert.Equal(t, "TEST_VALUE", os.Getenv("TEST_KEY"))
}

func TestLoad_WithOverloadParam(t *testing.T) {
	_ = godotenv.Write(map[string]string{"TEST_KEY": "NEW_VALUE"}, ".test.env")
	defer os.Remove(".test.env")

	os.Setenv("TEST_KEY", "OLD_VALUE")

	udotEnv := &udotEnvType{
		EnvParam:      stringSlice{".test.env"},
		OverloadParam: true,
	}

	assert.NotPanics(t, func() {
		udotEnv.Load()
	})

	assert.Equal(t, "NEW_VALUE", os.Getenv("TEST_KEY"))
}
