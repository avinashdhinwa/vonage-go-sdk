package nexmo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type verificationResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	ErrorText string `json:"error_text"`
}

func (response verificationResponse) StatusCode() (int, error) {
	return strconv.Atoi(response.Status)
}

func (client nexmoClient) Verify(number, brand, from string, length int, locale string) (string, error) {
	// TODO: Missing a couple of params, and locale is unused.
	// TODO: Timeouts are currently ignored.

	if length > 0 && length != 4 && length != 6 {
		return "", fmt.Errorf("code length must be 4 or 6")
	}

	params := url.Values{}
	params.Set("api_key", client.apiKey)
	params.Set("api_secret", client.apiSecret)
	params.Set("number", number)
	params.Set("brand", brand)

	url, err := url.Parse(client.baseURL + pathVerify)
	if err != nil {
		return "", err
	}
	url.RawQuery = params.Encode()

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Unexpected response from server: %d %s", response.StatusCode, response.Status)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return parseVerifyResponse(bytes)
}

func parseVerifyResponse(data []byte) (string, error) {
	response := verificationResponse{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	statusCode, err := response.StatusCode()
	if err != nil {
		return "", err
	}
	if statusCode != 0 {
		return "", fmt.Errorf("%s: %s", response.Status, response.ErrorText)
	}

	return response.RequestID, err
}