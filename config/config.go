package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"test-assignment-cookie-sync/cipher"
	"time"
)

const (
	DEFAULT_HOST = "0.0.0.0"
	DEFAULT_PORT = 9000
)

type Config struct {
	HTTPConfig
	StorageConfig
}

func (cfg *Config) validate() error {
	return buildError(map[string][]error{
		"HTTPConfig":    cfg.HTTPConfig.validate(),
		"StorageConfig": cfg.StorageConfig.validate(),
	})
}

func buildError(errsm map[string][]error) error {
	if len(errsm) == 0 || isEmpty(errsm) {
		return nil
	}
	b, err := json.Marshal(errsm)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return errors.New(string(b))
}

func isEmpty(errsm map[string][]error) bool {
	var counter int
	for _, v := range errsm {
		if len(v) != 0 {
			counter++
		}
	}
	return counter == 0
}

// ParseFlags parses specified flags and return error if some error occured
func ParseFlags() (*Config, error) {
	var config Config

	flag.StringVar(&config.HTTPConfig.Host, "host", DEFAULT_HOST, "flag for setup host of sync")
	flag.IntVar(&config.HTTPConfig.Port, "port", DEFAULT_PORT, "flag for setup port of sync")
	flag.StringVar(&config.HTTPConfig.Timeout, "timeout", "", "flag for setup timeout")

	flag.StringVar(&config.StorageConfig.Driver, "db-driver", "", "database driver")
	flag.StringVar(&config.StorageConfig.Host, "db-host", "", "database host")
	flag.IntVar(&config.StorageConfig.Port, "db-port", 0, "database port ")
	flag.BoolVar(&config.StorageConfig.IsEncrypted, "is-encrypted", false, "flag to enable decryption")
	flag.StringVar(&config.StorageConfig.User, "db-user", "", "database username")
	flag.StringVar(&config.StorageConfig.Password, "db-password", "", "database password")
	flag.StringVar(&config.StorageConfig.Database, "db-name", "", "database name")
	flag.StringVar(&config.StorageConfig.Schema, "db-schema", "", "database schema")
	flag.StringVar(&config.StorageConfig.Sslmode, "db-mode", "", "database mode")
	flag.BoolVar(&config.StorageConfig.IsDsnEnabled, "is-dsn-enabled", false, "flag to enable dsn")
	flag.StringVar(&config.StorageConfig.Dsn, "dsn", "", "already configured full path to db")

	flag.Parse()

	if config.StorageConfig.IsEncrypted {
		var err error
		config.StorageConfig.User, err = cipher.DecryptAES([]byte(config.StorageConfig.User), []byte(os.Getenv("SEC_SECRET")))
		if err != nil {
			return nil, fmt.Errorf("decrypt AES: %v", err)
		}
		config.StorageConfig.Password, err = cipher.DecryptAES([]byte(config.StorageConfig.Password), []byte(os.Getenv("SEC_SECRET")))
		if err != nil {
			return nil, fmt.Errorf("decrypt AES: %v", err)
		}
	}

	return &config, config.validate()
}

type HTTPConfig struct {
	Host    string
	Port    int
	Timeout string
}

func (h *HTTPConfig) validate() []error {
	var errs = make([]error, 0)
	if h.Host == "" {
		errs = append(errs, errors.New("host field is empty"))
	}
	if h.Port <= 0 {
		errs = append(errs, errors.New("port field is empty"))
	}
	_, err := time.ParseDuration(h.Timeout)
	if err != nil {
		errs = append(errs, errors.New("timeout field is invalid"))
	}
	return errs
}

type StorageConfig struct {
	Driver       string
	User         string
	Password     string
	Database     string
	Schema       string
	Sslmode      string
	Dsn          string
	Host         string
	Port         int
	IsDsnEnabled bool
	IsEncrypted  bool
}

func (st *StorageConfig) validate() []error {
	if st.IsDsnEnabled && st.Dsn != "" {
		return nil
	}

	var errs = make([]error, 0)

	if st.Driver == "" {
		errs = append(errs, errors.New("driver must not be empty"))
	}
	if st.Host == "" {
		errs = append(errs, errors.New("host must not be empty"))
	}
	if st.Port <= 0 {
		errs = append(errs, errors.New("port must not be negative"))
	}
	if st.User == "" {
		errs = append(errs, errors.New("user must not be empty"))
	}
	if st.Password == "" {
		errs = append(errs, errors.New("password must not be empty"))
	}
	if st.Database == "" {
		errs = append(errs, errors.New("database must not be empty"))
	}
	if st.Schema == "" {
		errs = append(errs, errors.New("schema must not be empty"))
	}
	if st.Sslmode == "" {
		errs = append(errs, errors.New("sslmode must not be empty"))
	}
	return errs
}
