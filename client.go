package uaago

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
    "crypto/tls"
)

type Client struct {
	uaaUrl string
}

func NewClient(uaaUrl string) Client {
	return Client{
		uaaUrl: uaaUrl,
	}
}

func (client Client) GetAuthToken(username, password string, insecureSkipVerify bool ) (string, error) {
	data := url.Values{"client_id": {username}, "grant_type": {"client_credentials"}}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth/token", client.uaaUrl), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(username, password)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    config := &tls.Config{InsecureSkipVerify: insecureSkipVerify}
    tr := &http.Transport{ TLSClientConfig: config }
    httpClient := &http.Client{Transport: tr}


    resp, err := httpClient.Do(request)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received a status code %v", resp.Status)
	}

	jsonData := make(map[string]interface{})
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jsonData)

	return fmt.Sprintf("%s %s", jsonData["token_type"], jsonData["access_token"]), err
}
