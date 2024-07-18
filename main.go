package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type BwgoConfig struct {
	cacheDir string
}

type BitwardenInstance struct {
	apiUrl       url.URL
	identityUrl  url.URL
	bearerToken  string
	clientId     string
	clientSecret string
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
		},
		bearerToken: "",
	}

	if err := bwi.oauth2(); err != nil {
		fmt.Println("OAUTH seems to have failed: %s", err.Error())
		return
	}

}

func (i *BitwardenInstance) oauth2() error {
	reader := strings.NewReader(fmt.Sprintf("grant_type=client_credential&client_id=%s&client_secret=%s", i.clientId, i.clientSecret))
	rc := io.NopCloser(reader)
	req := http.Request{Method: "POST", URL: &i.identityUrl, Proto: "http/2.0", Body: rc}
	resp, err := http.Post(i.identityUrl, "application/x-www-from-urlencoded", reader)
	if err != nil {
		return fmt.Errorf("get of API failed: %s", err.Error())
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("get of API failed: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received invalid status code back: %d", resp.StatusCode)
	}
	fmt.Println(string(content))
	fmt.Println(len(content))
	return nil
}

func getPassword(bi BitwardenInstance, name string) error {
	req := http.Request{Method: "", URL: bi.apiUrl}
}
