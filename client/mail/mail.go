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
	host     string
	port     int
	username string
	password string
	identity string
}

func New(cfg *Config) *Client {
	return &Client{
		host     : cfg.Host,
		port     : cfg.Port,
		username : cfg.Username,
		password : cfg.Password,
		identity : cfg.Identity,
	}
}

func (c *Client) Send(to string, message string, subject string) error {
	auth := smtp.PlainAuth(c.identity, c.username, c.password, c.host)

	msg := []byte(
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n",
	)

	return smtp.SendMail(fmt.Sprintf("%s:%d", c.host, c.port), auth, c.username, []string{to}, msg)
}
