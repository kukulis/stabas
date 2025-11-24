# stabas
valdymo konsolė

# Projekto paleidimas
Aplanke *./application* vykdyti komandą:

`go run .`


# sertifikatų generavimas

Eiti į tls direktoriją cmdline.

    openssl genrsa -out server.key 2048
    
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

# Naudojamos bibliotekos

* github.com/gin-gonic/gin – HTTP tinklo karkasas.
