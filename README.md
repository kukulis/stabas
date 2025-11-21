# stabas
valdymo konsolė

# paleisti projektą
Aplanke *./application* vykdyti komandą:

`go run .`


# sertifikatų generavimas

Eiti į tls direktoriją cmdline.

    openssl genrsa -out server.key 2048
    
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650


# run single test

    go test -run TestGroupTasks


# build for windows

    GOOS=windows GOARCH=amd64 go build