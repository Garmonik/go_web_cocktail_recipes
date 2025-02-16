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
	SecretKey  string `yaml:"secret_key" env-required:"true"`
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
	fmt.Println("111111111111111111111111111111111111111111111111111111111111111111111")
	privBytes, _ := x509.MarshalECPrivateKey(cfg.JWT.PrivateKey)

	fmt.Println(string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})))

	// Публичный ключ
	pubBytes, _ := x509.MarshalPKIXPublicKey(cfg.JWT.PublicKey)
	fmt.Println(string(pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: pubBytes})))
	fmt.Println("111111111111111111111111111111111111111111111111111111111111111111111")
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
	// Парсим приватный ключ
	privKey, err := parseECDSAKey(cfg.JWT.PrivateKeyPEM, true)
	if err != nil {
		return fmt.Errorf("ошибка загрузки приватного ключа: %w", err)
	}

	privateKey, ok := privKey.(*ecdsa.PrivateKey)
	if !ok {
		return fmt.Errorf("неверный формат приватного ключа")
	}
	cfg.JWT.PrivateKey = privateKey // ✅ Теперь это *ecdsa.PrivateKey

	// Парсим публичный ключ
	pubKey, err := parseECDSAKey(cfg.JWT.PublicKeyPEM, false)
	if err != nil {
		return fmt.Errorf("ошибка загрузки публичного ключа: %w", err)
	}

	publicKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("неверный формат публичного ключа")
	}
	cfg.JWT.PublicKey = publicKey // ✅ Теперь это *ecdsa.PublicKey

	return nil
}
