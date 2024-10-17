package mail

import (
	"fmt"
	"net/smtp"
)

// Config holds the configuration for the email client.
// It includes fields for the mail server host, port, username, password, and identity.
type Config struct {
	Host     string `yaml:"host"     env:"MAIL_HOST"`     // The SMTP server hostname
	Port     int    `yaml:"port"     env:"MAIL_PORT"`     // The SMTP server port
	Username string `yaml:"username" env:"MAIL_USERNAME"` // The username for authentication
	Password string `env:"MAIL_PASSWORD" `                // The password for authentication
	Identity string `yaml:"identity" env:"MAIL_IDENTITY"` // The identity of the sender
}

// Client represents an email client that can send emails using the specified configuration.
type Client struct {
	host     string
	port     int
	username string
	password string
	identity string
}

// New creates a new Client instance using the provided configuration.
// It initializes the Client with the host, port, username, password, and identity.
func New(cfg *Config) *Client {
	return &Client{
		host:     cfg.Host,
		port:     cfg.Port,
		username: cfg.Username,
		password: cfg.Password,
		identity: cfg.Identity,
	}
}

// Send sends an email to the specified recipient with the given message, subject, and content type.
// It uses SMTP authentication and constructs the email message in the required format.
func (c *Client) Send(to string, message string, subject string, contentType string) error {
	auth := smtp.PlainAuth(c.identity, c.username, c.password, c.host)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\n%s\n\n%s\r\n",
		c.username, to, subject, contentType, message,
	))

	return smtp.SendMail(fmt.Sprintf("%s:%d", c.host, c.port), auth, c.username, []string{to}, msg)
}

// Mailing sends a message to a list of recipients with the given message, subject, and content type.
func (c *Client) Mailing(emails []string, message string, subject string, contentType string) error {
	for i := 0; i < len(emails); i++ {
		err := c.Send(emails[i], message, subject, contentType)
		if err != nil {
			return err
		}
	}

	return nil
}

// PersonalMailing sends a message to a list of recipients substituting values into given template.
// The template format should be the same as for fmt`s functions.
// The values are taken from a multiple, where the key is the recipient's mail, and the value is an array of values
// that will be substituted into the template. Be careful with the order of values in the template and the array of values.
func (c *Client) PersonalMailing(values map[string][]interface{}, template string, subject string, contentType string) error {
	for email, vls := range values {
		message := fmt.Sprintf(template, vls...)

		err := c.Send(email, message, subject, contentType)
		if err != nil {
			return err
		}
	}

	return nil
}
