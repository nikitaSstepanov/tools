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
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		identity: cfg.Identity,
	}
}

func (c *Client) Send(to string, message string, subject string, contentType string) error {
	auth := smtp.PlainAuth(c.identity, c.username, c.password, c.host)

	msg := []byte(fmt.Sprintf(
		`
			From: %s \r\n 
			To: %s \r\n 
			Subject: %s \n
			%s \r\n \r\n 
			%s \r\n
		`,
		c.username, to, subject, contentType, message,
	))

	return smtp.SendMail(fmt.Sprintf("%s:%d", c.host, c.port), auth, c.username, []string{to}, msg)
}
