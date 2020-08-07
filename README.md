# api-email

### Run
1. Run go get github.com/lizzle-shizzle/api-email/api
2. Run chmod u+x build.sh in project directory
3. Run ./build.sh
4. Run docker-compose up

### Endpoints
- POST localhost:8080/email
	- Request body:
	```
	{
		"subject": "Testing",
		"body": "<b>Hello world</b>",
		"email": "liz@example.com"
	}
	```
	- Authorization: Bearer [token sent in email]
