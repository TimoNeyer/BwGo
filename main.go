package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BitwardenInstance struct {
	apiUrl         url.URL
	identityUrl    url.URL
	preauthUrl     url.URL
	bearerToken    string
	clientId       string
	clientSecret   string
	SessionTimeout int
}

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	bwi := BitwardenInstance{
		apiUrl: url.URL{
			Scheme: "https",
			Opaque: "api.bitwarden.com",
		},
		identityUrl: url.URL{
			Scheme: "https",
			Opaque: "identity.bitwarden.com",
			Path:   "/connect/token",
		},
		preauthUrl: url.URL{
			Scheme: "https",
			Opaque: "vault.bitwarden.com",
		},
		bearerToken: "",
	}

	if err := bwi.oauth2(); err != nil {
		fmt.Println("OAUTH seems to have failed: %s", err.Error())
		return
	}

}

func (i *BitwardenInstance) preauth() error {
	return nil
}

func (i *BitwardenInstance) oauth2() error {
	reader := strings.NewReader(fmt.Sprintf("grant_type=client_secret&client_id=%s&client_secret=%s", i.clientId, i.clientSecret))
	rc := io.NopCloser(reader)
	req := http.Request{
		Method: "POST",
		URL:    &i.identityUrl,
		Proto:  "http/2.0",
		Body:   rc,
	}
	req.Header.Set("Content-Type", "application/z-www-form-urlencoded")
	client := http.Client{}
	respone, err := client.Do(&req)
	if err != nil {
		return fmt.Errorf("failed to form request: %s", err.Error())
	}
	if respone.StatusCode != http.StatusOK {
		if respone.StatusCode == http.StatusForbidden {
			return fmt.Errorf("Error: Access denied (received Status Forbidden: %d)", respone.StatusCode)
		} else {
			return fmt.Errorf("Error: Received invalid Status code: %d", respone.StatusCode)
		}
	}
	defer respone.Body.Close()
	//	resp, err := http.Post(i.identityUrl, "application/x-www-from-urlencoded", reader)
	body, err := io.ReadAll(respone.Body)
	if err != nil {
		return fmt.Errorf("get of API failed: %s", err.Error())
	}
	var token Token
	if err := json.Unmarshal(body, token); err != nil {
		return fmt.Errorf("Error: failed to unpack token: %s", err.Error())
	}
	i.bearerToken = token.AccessToken
	i.SessionTimeout = token.ExpiresIn
	return nil
}
