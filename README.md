# stabas
valdymo konsolė

# Projekto paleidimas

## Susigeneruoti sertifikatus 

Žr. sekciją "Sertifikatų generavimas" dokumentacijoje.

## Paleisti aplikaciją

Aplanke *./application* vykdyti komandą:

`go run .`


# sertifikatų generavimas

Eiti į tls direktoriją cmdline.

    openssl genrsa -out server.key 2048
    
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

Be prompto:

Pavyzdys:

(openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname")

    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -nodes -subj "/C=LT/ST=Lietuva/L=Kaunas/O=Darbelis/OU=stabas/CN=stabas"



# run single test

    go test -run TestGroupTasks
# Naudojamos bibliotekos

* github.com/gin-gonic/gin – HTTP tinklo karkasas.

# build for windows

    GOOS=windows GOARCH=amd64 go build

# run all tests

    cd application

    go test ./...

