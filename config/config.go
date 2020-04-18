package config

import (
	"bytes"
	"os"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/viper"
	"github.com/vicanso/alarms/validate"
)

var (
	box     = packr.New("config", "../configs")
	env     = os.Getenv("GO_ENV")
	appName string
)

type (
	// MailConfig mail's config
	MailConfig struct {
		Host     string `valid:"host"`
		Port     int    `valid:"port"`
		User     string `valid:"email"`
		Password string `valid:"runelength(1|100)"`
	}
)

const (
	// Dev development env
	Dev = "dev"
	// Test test env
	Test = "test"
	// Production production env
	Production = "production"

	defaultListen = ":7001"
)

var (
	defaultViper = viper.New()
)

func init() {
	configType := "yml"
	configExt := "." + configType
	data, err := box.Find("default" + configExt)
	if err != nil {
		panic(err)
	}
	viper.SetConfigType(configType)
	v := viper.New()
	v.SetConfigType(configType)
	// 读取默认配置中的所有配置
	err = v.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	configs := v.AllSettings()
	// 将default中的配置全部以默认配置写入
	for k, v := range configs {
		defaultViper.SetDefault(k, v)
	}

	// 根据当前运行环境配置读取
	envConfigFile := GetENV() + configExt
	data, err = box.Find(envConfigFile)
	if err != nil {
		panic(err)
	}
	// 读取当前运行环境对应的配置
	err = defaultViper.ReadConfig(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	appName = GetString("app")
}

func validatePanic(v interface{}) {
	err := validate.Do(v, nil)
	if err != nil {
		panic(err)
	}
}

func GetAppName() string {
	return appName
}

// GetENV get go env
func GetENV() string {
	if env == "" {
		return Dev
	}
	return env
}

// GetInt viper get int
func GetInt(key string) int {
	return defaultViper.GetInt(key)
}

// GetUint viper get uint
func GetUint(key string) uint {
	return defaultViper.GetUint(key)
}

// GetUint32 viper get uint32
func GetUint32(key string) uint32 {
	return defaultViper.GetUint32(key)
}

// GetIntDefault get int with default value
func GetIntDefault(key string, defaultValue int) int {
	v := GetInt(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetUint32Default get uint32 with default value
func GetUint32Default(key string, defaultValue uint32) uint32 {
	v := GetUint32(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetString viper get string
func GetString(key string) string {
	return defaultViper.GetString(key)
}

// GetStringDefault get string with default value
func GetStringDefault(key, defaultValue string) string {
	v := GetString(key)
	if v != "" {
		return v
	}
	return defaultValue
}

// GetDuration viper get duration
func GetDuration(key string) time.Duration {
	return defaultViper.GetDuration(key)
}

// GetDurationDefault get duration with default value
func GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	v := GetDuration(key)
	if v != 0 {
		return v
	}
	return defaultValue
}

// GetStringSlice viper get string slice
func GetStringSlice(key string) []string {
	return defaultViper.GetStringSlice(key)
}

// GetStringMap get string map
func GetStringMap(key string) map[string]interface{} {
	return defaultViper.GetStringMap(key)
}

// GetListen get server listen address
func GetListen() string {
	return GetStringDefault("listen", defaultListen)
}

// GetMailConfig get mail config
func GetMailConfig() MailConfig {
	prefix := "mail."
	pass := GetString(prefix + "password")
	if os.Getenv(pass) != "" {
		pass = os.Getenv(pass)
	}
	mailConfig := MailConfig{
		Host:     GetString(prefix + "host"),
		Port:     GetInt(prefix + "port"),
		User:     GetString(prefix + "user"),
		Password: pass,
	}
	validatePanic(&mailConfig)
	return mailConfig
}
