# api-email

### Run
1. Run ./build.sh
2. Run docker-compose up

### Endpoints
- POST localhost:8080/email
	- Request body:
	{
		"subject": "Testing",
		"body": "<b>Hello world</b>",
		"email": "liz@example.com"
	}
	- Authorization: Bearer <<token sent in email>>