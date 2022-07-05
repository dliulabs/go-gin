# Seed users database

```
cd users
go run main.go
```

# Test run

```
go run ./cmd/api
```

```
curl -X POST http://localhost:8080/signin -d '{"usern ame":"david","password":"1234"}'

curl -X POST http://localhost:8080/signin -d '{"username":"admin","password":"password"}' |  jq -r

curl -X POST http://localhost:8080/signin -d '{"username":"packt","password":"RE4zfHB35VPtTkbT"}' |  jq -r

curl -X POST http://localhost:8080/signin -d '{"username":"mlabouardy","password":"L3nSFRcZzNQ67bcc"}' |  jq -r

```

```
curl -b cookie.txt --cookie-jar cookie.txt -X POST http://localhost:8080/signin -d '{"username":"mlabouardy","password":"L3nSFRcZzNQ67bcc"}' | export TOKEN=`jq '.token' | tr -d '"'`

curl -b cookie.txt --cookie-jar cookie.txt -X GET http://localhost:8080/recipes -H "Authorization: Bearer ${TOKEN}"
```

```
curl  -X POST http://localhost:8080/refresh -H "Authorization: ${TOKEN}"
```