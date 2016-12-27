FROM alpine:3.4
MAINTAINER Roberto J. Rojas <robertojrojas@gmail.com>

ADD dogs-service /usr/bin/dogs-service

ENV RABBITMQ_URI amqp://guest:guest@localhost:5672/
ENV RABBITMQ_QUEUE dog_service_rpc_queue
ENV MONGODB_URI mongodb://localhost:27017
ENV MONGODB_DBNAME dogsDB
ENV MONGODB_COLLECTION dogs

ENTRYPOINT ["dogs-service"]