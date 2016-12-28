kubectl create configmap rabbitmq-config --from-literal=rabbitmq.uri=amqp://guest:guest@10.0.1.122:5673/ --from-literal=rabbitmq.queue=dog_service_rpc_queue
kubectl create configmap mongodb-config --from-literal=mongodb.uri=mongodb://10.0.1.122:27017 --from-literal=mongodb.dbname=dogsDB --from-literal=mongodb.collection=dogs
kubectl create -f dogs-service-pod.yaml
