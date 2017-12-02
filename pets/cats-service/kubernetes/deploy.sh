#!/bin/bash

mysql_user=${1:-catuser}
mysql_password=${2:-my-secret-pw}
mysql_cert=$3

echo -n "${mysql_user}" > ./username.txt
echo -n "${mysql_password}" > ./password.txt

kubectl create secret generic mysql-secrets \
      --from-file=username=./username.txt \
      --from-file=password=./password.txt \
      --dry-run -o yaml > mysql-secrets.yml
rm -f ./username.txt
rm -f ./password.txt
kubectl create -f mysql-secrets.yml

if [[ -n ${mysql_cert} ]];
then
  kubectl create secret generic mysql-secrets \
      --from-file=mysql-cert=${mysql_cert}    \
      --dry-run -o yaml | kubectl apply -f -
fi

kubectl create -f mysql-endpoints.yml
kubectl create -f mysql-svc.yml
kubectl create -f mysql-configmap.yml
kubectl create -f cats-service-pod.yml
kubectl create -f cats-service-svc.yml
