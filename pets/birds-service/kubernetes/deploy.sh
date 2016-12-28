kubectl create configmap postgresql-config --from-literal=postgresql.uri="user=birdman password=mycape dbname=birddb sslmode=disable"
kubectl create -f birds-service-pod.yaml
