# Web-forum

## Project Description
Web-Forum is a website similar to Reddit that allows users to view and create posts related to a specific topic. The project is done using Golang, Javascript, HTML, and CSS. The backend is fully done using Go and uses sqlite3 to store data.

## Functionalities
- User Authentication: Sign-in and Sign-up
- User Authorization: Session Management
- Post Creation
- Post Filtering
- Comments
- Like and Dislike Posts and Comments


## Usage
- Make sure to have Golang installed
- Clone the repository `https://github.com/Mahmood-7sain/web-forum.git`
- Navigate to the correct directory: `cd forum/forum`
- Get the external packages: `go get github.com/mattn/go-sqlite3` and `go get golang.org/x/crypto`
- Run the command `go mod tidy`

### Running locally
- Run the command `go run main.go`
- Navigate to `localhost:8080` and use the website

### Running on Docker
If you have docker installed then you can use the ready-made bash file to run the website on docker containers
- Run the command `bash buildAndRun.sh` and follow the prompts
- The container will run automatically on the port you specified
- To remove the images and containers use the `bash prune.sh` 



