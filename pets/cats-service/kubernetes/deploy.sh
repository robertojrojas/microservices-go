kubectl create configmap mysql-config --from-literal=mysql.uri=root:my-secret-pw@tcp(:3306)/cats_db
kubectl create -f cats-service-pod.yaml
