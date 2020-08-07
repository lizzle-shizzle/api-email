package api

import (
	"time"
	"fmt"
)

type Email struct {
	Subject string `json:"subject"`
	Body string `json:"body"`
	Email string `json:"email"`
}

type EmailLog struct {
	ID int `json:"id"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	Email string `json:"email"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
}

// AddNewEmailLog adds a new email log to the email log table
func AddNewEmailLog(c *Client, email Email) (EmailLog, error) {
	currentTimestamp := time.Now()
	sqlStatement := `
	INSERT INTO email_log (subject, body, email, created_timestamp) VALUES ($1, $2, $3, $4)
	RETURNING id`

	id := 0
	err := c.DB.QueryRow(sqlStatement, email.Subject, email.Body, email.Email, currentTimestamp).Scan(&id)
	if err != nil {
		fmt.Errorf("Failed to insert email log: %s", err)
		return EmailLog{}, err
	}

	return EmailLog{
		id,
		email.Subject,
		email.Body,
		email.Email,
		currentTimestamp,
	}, nil
}
