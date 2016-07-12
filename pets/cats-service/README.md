
docker run --name some-mysql -v `pwd`/datadir:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql

curl -v -H "Content-Type: application/json" -X POST -d '{"name":"second","age":1,"type":"fatty"}' http://localhost:8091/api/cats

curl -v http://localhost:8091/api/cats

curl -v -H "Content-Type: application/json" -X PUT -d '{"name":"second","age":1,"type":"biggie"}' http://localhost:8091/api/cats/2

curl -v -X DELETE http://localhost:8091/api/cats/2
