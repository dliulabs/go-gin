# Auth0 Info

```
curl --request POST \
  --url https://dev-osud0o9f.auth0.com/oauth/token \
  --header 'content-type: application/json' \
  --data '{"client_id":"clientid","client_secret":"secret","audience":"https://api.recipes.io","grant_type":"client_credentials"}'
```

profile: https://dev-osud0o9f.auth0.com/.well-known/jwks.json

```
{
  "iss": "https://dev-osud0o9f.auth0.com/",
  "sub": "q1OEUARa2q6QsQrRq0J8yDl2lyAfGfDb@clients",
  "aud": "https://api.recipes.io",
  "iat": 1656943851,
  "exp": 1657030251,
  "azp": "q1OEUARa2q6QsQrRq0J8yDl2lyAfGfDb",
  "gty": "client-credentials"
}
```

# Test Run

```
curl  -X GET http://localhost:8080/recipes | jq -r
curl  -X GET http://localhost:8080/recipes/62bfa70fa01938df3d3eb76a | jq -r

go run ./auth0-client/main.go | export TOKEN=`jq '.access_token' | tr -d '"'`
curl -H "Authorization: Bearer ${TOKEN}" -X GET http://localhost:8080/recipes/62bfa70fa01938df3d3eb76a | jq -r
```

```
ngrok http 8080

go run ./auth0-client/main.go | export TOKEN=`jq '.access_token' | tr -d '"'`

curl -H "Authorization: Bearer ${TOKEN}" -X GET https://903b-2601-484-c500-73a-bc8b-32-2973-670f.ngrok.io/recipes/62c30612011537e9d8f2774c | jq -r
```

```
mkdir certs
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/localhost.key -out certs/localhost.crt -config ./certs/server.conf

curl --cacert certs/localhost.crt https://api.recipes.io/recipes --insecure
```