package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"bytes"
	"strings"
)

type App struct {
	Client *Client
}

type Client struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	DB       *sql.DB
}

func (c *Client) connect() (*sql.DB, error) {
	params := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)

	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	return db, err
}

func (a *App) Initialize() error {
	var err error
	database, err := a.Client.connect()
	if err != nil {
		log.Println("Failed to connect to database")
		return err
	}
	a.Client.DB = database

	// Initialize route
	http.HandleFunc("/email", a.SendEmail)

	return nil
}

func (a *App) Run(addr string) {
	log.Printf(fmt.Sprintf("Email API listening on %v...\n", addr[1:]))
	log.Fatalln(http.ListenAndServe(addr, nil))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal response body: %s\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) SendEmail(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		var email Email

		decoder := json.NewDecoder(r.Body)
		if r.Body == nil {
			log.Println("Send email failed - Invalid payload request: empty body")
			respondWithError(w, http.StatusBadRequest, "Invalid payload request: empty body")
			return
		}
		if err := decoder.Decode(&email); err != nil {
			log.Printf("Failed to decode message: %s\n", err)
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()

		requestBody := []byte(
			`{
				"personalizations": [
					{
						"to": [
							{
								"email": "` + email.Email + `"
							}
						],
						"subject": "` + email.Subject + `"
					}
				],
				"from": {
					"email": "elle.bee.elle.bee@gmail.com",
					"name": "Liz"
				},
				"content": [
					{
						"type": "text/html",
						"value": "` + email.Body + `"
					}
				]
			}`)

		client := &http.Client{}

		req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer " + reqToken)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == 202 {
			log.Printf("Email sent successfully: Status - %v\n", resp.StatusCode)
		} else {
			log.Fatalf("Failed to send email: Status - %v\n", resp.StatusCode)
		}

		msg, err := AddNewEmailLog(a.Client, email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, msg)
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}