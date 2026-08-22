package config

import (
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	masterKeySize                 = 32
	minimumPaginationSecretLength = 32
	defaultPaginationTTL          = 15 * time.Minute
	defaultDatabaseMaxOpen        = 10
	defaultDatabaseMaxIdle        = 5
	defaultDatabaseConnLifetime   = 30 * time.Minute
)

var (
	ErrInvalidRuntimeConfiguration = errors.New("invalid Akritas runtime configuration")
	ErrInvalidMasterKey            = errors.New("AKRITAS_MASTER_KEY must be valid Base64 encoding exactly 32 bytes")
)

// Config is the single validated runtime configuration consumed by bootstrap.
// Secret values are decoded during loading so adapters never need to parse env.
type Config struct {
	DatabaseURL                   string        `mapstructure:"AKRITAS_DATABASE_URL"`
	PublicURL                     string        `mapstructure:"AKRITAS_PUBLIC_URL"`
	MasterKeyEncoded              string        `mapstructure:"AKRITAS_MASTER_KEY"`
	PaginationSecretValue         string        `mapstructure:"AKRITAS_PAGINATION_SECRET"`
	PaginationTTL                 time.Duration `mapstructure:"AKRITAS_PAGINATION_TTL"`
	DatabaseMaxOpenConnections    int           `mapstructure:"AKRITAS_DB_MAX_OPEN_CONNECTIONS"`
	DatabaseMaxIdleConnections    int           `mapstructure:"AKRITAS_DB_MAX_IDLE_CONNECTIONS"`
	DatabaseConnectionMaxLifetime time.Duration `mapstructure:"AKRITAS_DB_CONNECTION_MAX_LIFETIME"`

	// Derived secrets are populated only after validation. The raw fields above
	// are cleared before Config leaves this package.
	MasterKey        []byte `mapstructure:"-"`
	PaginationSecret []byte `mapstructure:"-"`
}

func Load() (Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.SetConfigName("app")
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindEnvironment(v)
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) && !os.IsNotExist(err) {
			return Config{}, ErrInvalidRuntimeConfiguration
		}
	}
	return loadFromViper(v)
}

// LoadValues is retained as a deterministic test seam while still exercising
// the same Viper defaults, unmarshalling and validation path as Load.
func LoadValues(getenv func(string) string) (Config, error) {
	if getenv == nil {
		return Config{}, ErrInvalidRuntimeConfiguration
	}
	v := viper.New()
	for _, key := range configurationKeys() {
		if value := getenv(key); value != "" {
			v.Set(key, value)
		}
	}
	return loadFromViper(v)
}

func loadFromViper(v *viper.Viper) (Config, error) {
	if v == nil {
		return Config{}, ErrInvalidRuntimeConfiguration
	}
	setDefaults(v)
	var raw Config
	if err := v.Unmarshal(&raw); err != nil {
		return Config{}, ErrInvalidRuntimeConfiguration
	}

	databaseURL := strings.TrimSpace(raw.DatabaseURL)
	parsedDatabase, err := url.Parse(databaseURL)
	if err != nil || (parsedDatabase.Scheme != "postgres" && parsedDatabase.Scheme != "postgresql") || parsedDatabase.Host == "" {
		return Config{}, ErrInvalidRuntimeConfiguration
	}
	publicURL := strings.TrimRight(strings.TrimSpace(raw.PublicURL), "/")
	parsedPublic, err := url.Parse(publicURL)
	if err != nil || parsedPublic.Scheme != "https" || parsedPublic.Host == "" || parsedPublic.User != nil || parsedPublic.Path != "" || parsedPublic.RawQuery != "" || parsedPublic.Fragment != "" {
		return Config{}, ErrInvalidRuntimeConfiguration
	}
	masterKey, err := ParseMasterKey(strings.TrimSpace(raw.MasterKeyEncoded))
	if err != nil {
		return Config{}, err
	}
	paginationSecret := []byte(raw.PaginationSecretValue)
	if len(paginationSecret) < minimumPaginationSecretLength || raw.PaginationTTL <= 0 || raw.DatabaseMaxOpenConnections <= 0 || raw.DatabaseMaxIdleConnections < 0 || raw.DatabaseMaxIdleConnections > raw.DatabaseMaxOpenConnections || raw.DatabaseConnectionMaxLifetime <= 0 {
		clear(masterKey)
		return Config{}, ErrInvalidRuntimeConfiguration
	}
	raw.DatabaseURL = databaseURL
	raw.PublicURL = publicURL
	raw.MasterKeyEncoded = ""
	raw.PaginationSecretValue = ""
	raw.MasterKey = masterKey
	raw.PaginationSecret = append([]byte(nil), paginationSecret...)
	return raw, nil
}

func ParseMasterKey(value string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(decoded) != masterKeySize {
		return nil, ErrInvalidMasterKey
	}
	key := make([]byte, len(decoded))
	copy(key, decoded)
	return key, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("AKRITAS_PAGINATION_TTL", defaultPaginationTTL)
	v.SetDefault("AKRITAS_DB_MAX_OPEN_CONNECTIONS", defaultDatabaseMaxOpen)
	v.SetDefault("AKRITAS_DB_MAX_IDLE_CONNECTIONS", defaultDatabaseMaxIdle)
	v.SetDefault("AKRITAS_DB_CONNECTION_MAX_LIFETIME", defaultDatabaseConnLifetime)
}

func bindEnvironment(v *viper.Viper) {
	for _, key := range configurationKeys() {
		_ = v.BindEnv(key)
	}
}

func configurationKeys() []string {
	typeOf := reflect.TypeOf(Config{})
	keys := make([]string, 0, typeOf.NumField())
	for i := 0; i < typeOf.NumField(); i++ {
		if key := typeOf.Field(i).Tag.Get("mapstructure"); key != "" && key != "-" {
			keys = append(keys, key)
		}
	}
	return keys
}
