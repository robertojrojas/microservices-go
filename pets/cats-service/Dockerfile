FROM alpine:3.4
MAINTAINER Roberto J. Rojas <robertojrojas@gmail.com>

ADD cats-service /usr/bin/cats-service

ENV MYSQL_DB_URI root:my-secret-pw@tcp(:3306)/cats_db

EXPOSE 8091

ENTRYPOINT ["cats-service"]