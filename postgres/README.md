# starting postgre

```
docker-compose -f postgres.yml up -d
```

edit casdoor/conf/app.conf to have the following:

```
driverName = postgres
dataSourceName = "user=postgres password=postgres host=localhost port=5432 sslmode=disable dbname=casdoor"
dbName =
```