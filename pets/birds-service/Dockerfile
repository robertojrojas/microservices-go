FROM alpine:3.4
MAINTAINER Roberto J. Rojas <robertojrojas@gmail.com>

ADD birds-service /usr/bin/birds-service

ENV POSTGRESQL_DB_URI "user=birdman password=mycape dbname=birddb sslmode=disable"

EXPOSE 8092

ENTRYPOINT ["birds-service"]