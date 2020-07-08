package config

import (
	"log"
	"os"
)

type Config struct {
	Host   string
	Port   string
	Scheme string
	Domain string
	Secret string
	DSN    string
	SMTP   SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func Load() *Config {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatalf("HOST not provided")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT not provided")
	}
	scheme := os.Getenv("SCHEME")
	if scheme == "" {
		log.Println("SCHEME not provided, defaulting to https")
		scheme = "https"
	}
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatalf("DOMAIN not provided")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatalf("SECRET not provided")
	}
	dsn := os.Getenv("POSTGRES_CONNECTION")
	if dsn == "" {
		log.Fatalf("POSTGRES_CONNECTION not provided")
	}
	smtpHost := os.Getenv("SMTPHOST")
	if smtpHost == "" {
		log.Fatalf("SMTPHOST not provided")
	}
	smtpPort := os.Getenv("SMTPPORT")
	if smtpPort == "" {
		log.Fatalf("SMTPPORT not provided")
	}
	smtpUsername := os.Getenv("SMTPUSERNAME")
	if smtpUsername == "" {
		log.Fatalf("SMTPUSERNAME not provided")
	}
	smtpPassword := os.Getenv("SMTPPASSWORD")
	if smtpPassword == "" {
		log.Fatalf("SMTPPASSWORD not provided")
	}

	return &Config{
		Host:   host,
		Port:   port,
		Scheme: scheme,
		Domain: domain,
		Secret: secret,
		DSN:    dsn,
		SMTP: SMTPConfig{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUsername,
			Password: smtpPassword,
		},
	}
}
