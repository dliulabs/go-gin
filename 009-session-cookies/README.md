# Test run

```
go run ./cmd/api/session.go

curl -s -I http://localhost:8080/login --cookie-jar cookie.txt

curl -s -b cookie.txt http://localhost:8080/secret
```

```
curl http://localhost:8000/login -b cookie.txt --cookie-jar cookie.txt
```