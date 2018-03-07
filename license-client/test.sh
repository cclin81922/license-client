#!/bin/bash

# use local dev license-server
#docker run -dit -h localhost.localdomain --name httpd-2.4 -p 8080:80 -p 8443:443 -v "$PWD/data":/usr/local/apache2/htdocs/ my-httpd-2.4
export LICENSE_SERVER=https://localhost.localdomain:8443/wsgi
go run main.go

# use remote dev license-server
#export LICENSE_SERVER=https://192.168.240.56.xip.io/wsgi
#go run main.go
