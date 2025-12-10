# HeadQuarters
Headquarters control system.

Handling tasks in runtime.

# Project startup

## Application launch

In folder *./application* run command:

`go run .`


## Generate certificates

Cmd line go to tls directory and run command:

    openssl genrsa -out server.key 2048
    
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

Example to run without questions prompt:

    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -nodes -subj "/C=LT/ST=Lietuva/L=Kaunas/O=Darbelis/OU=stabas/CN=stabas"


# Developing hints

## Running single test

    go test -run TestGroupTasks

## run all tests

    cd application

    go test ./...

## Used libraries

* github.com/gin-gonic/gin – HTTP network carcass.

## build for windows

    GOOS=windows GOARCH=amd64 go build


