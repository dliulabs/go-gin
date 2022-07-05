package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	url := "https://dev-osud0o9f.auth0.com/oauth/token"

	// payload := strings.NewReader("{\"client_id\":\"0wlUEexUtqFPF0Zru1b9YCAD2rDuI3BS\",\"client_secret\":\"AyHLhNKPBrU70KafdJlNVJc1QJ5MhdSIfjc4xF2FyzqClJLm1-qTSL66HNftsHLq\",\"audience\":\"https://api. recipes.io\",\"grant_type\":\"client_credentials\"}")
	payload := strings.NewReader(
		fmt.Sprintf(`{
			"client_id":"%s",
			"client_secret":"%s",
			"audience":"https://api.recipes.io",
			"grant_type":"client_credentials"
		}`, os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")))

	// fmt.Println(payload)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	// fmt.Println(res)
	fmt.Println(string(body))
}
