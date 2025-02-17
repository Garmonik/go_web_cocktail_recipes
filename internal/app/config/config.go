package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	DBPath     string `yaml:"db_path" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	JWT        `yaml:"jwt"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JWT struct {
	PrivateKeyPEM string `yaml:"private_key" env-required:"true"`
	PublicKeyPEM  string `yaml:"public_key" env-required:"true"`
	PrivateKey    *ecdsa.PrivateKey
	PublicKey     *ecdsa.PublicKey
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	err := cfg.ParseKeys()
	if err != nil {
		log.Fatalf("cannot ParseKeys: %s", err)
	}
	return &cfg
}

func parseECDSAKey(base64Key string, isPrivate bool) (interface{}, error) {
	decoded, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, fmt.Errorf("error decode Base64: %w", err)
	}

	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, fmt.Errorf("error decode PEM")
	}

	if isPrivate {
		return x509.ParseECPrivateKey(block.Bytes)
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public_key: %w", err)
	}

	return pubKey.(*ecdsa.PublicKey), nil
}

func (cfg *Config) ParseKeys() error {
	privKey, _ := parseECDSAKey(cfg.JWT.PrivateKeyPEM, true)

	privateKey, ok := privKey.(*ecdsa.PrivateKey)
	if !ok {
		return fmt.Errorf("error privateKey")
	}
	cfg.JWT.PrivateKey = privateKey

	pubKey, _ := parseECDSAKey(cfg.JWT.PublicKeyPEM, false)

	publicKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error publicKey")
	}
	cfg.JWT.PublicKey = publicKey

	return nil
}
