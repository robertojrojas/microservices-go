kubectl create configmap postgresql-config --from-literal=postgresql.uri=postgres://birdman:mycape@10.0.1.122/pqgotest?sslmode=disable
kubectl create -f birds-service-pod.yaml
