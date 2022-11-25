package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ReCAPTCHA_VerifyURL string
	ReCAPTCHA_ServerKey string
}

// GoogleRecaptchaResponse ...
type GoogleRecaptchaResponse struct {
	Success            bool     `json:"success"`
	ChallengeTimestamp string   `json:"challenge_ts"`
	Hostname           string   `json:"hostname"`
	ErrorCodes         []string `json:"error-codes"`
}

// NewConfig create a new instance of configuration
func NewConfig() (*Config, error) {
	return &Config{
		ReCAPTCHA_VerifyURL: "",
		ReCAPTCHA_ServerKey: "",
	}, nil
}

func (o *Config) LoadConfig() error {
	// load .env file from given path
	// we keep it empty it will load .env from current directory

	rootPath := os.Getenv("HOME")
	err := godotenv.Load(rootPath + "/.env")

	if err != nil {
		return err
	}

	o.ReCAPTCHA_VerifyURL = os.Getenv("ReCAPTCHA_VerifyURL")
	o.ReCAPTCHA_ServerKey = os.Getenv("ReCAPTCHA_ServerKey")

	return nil
}
