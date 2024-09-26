package mail

import (
	"fmt"
	"net/smtp"
)

type Config struct {
	Host     string `yaml:"host"     env:"MAIL_HOST"`
	Port     int    `yaml:"port"     env:"MAIL_PORT"`
	Username string `yaml:"username" env:"MAIL_USERNAME"`
	Password string `env:"MAIL_PASSWORD"`
	Identity string `yaml:"identity" env:"MAIL_IDENTITY"`
}

type Client struct {
	Host     string
	Port     int
	Username string
	Password string
	Identity string
}

func New(cfg *Config) *Client {
	return &Client{
		Host     : cfg.Host,
		Port     : cfg.Port,
		Username : cfg.Username,
		Password : cfg.Password,
		Identity : cfg.Identity,
	}
}

func (c *Client) Send(to string, message string, subject string) error {
	auth := smtp.PlainAuth(c.Identity, c.Username, c.Password, c.Host)

	msg := []byte(
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n",
	)

	return smtp.SendMail(fmt.Sprintf("%s:%d", c.Host, c.Port), auth, c.Username, []string{to}, msg)
}
