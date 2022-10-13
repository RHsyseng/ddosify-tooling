package ddosify

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

//getEnv returns the value for a given Env Var
func getEnv(varName string, defaultValue string) string {
	if varValue, ok := os.LookupEnv(varName); ok {
		return varValue
	}
	return defaultValue
}

//ValidateURL validates that the URL is correct
func ValidateURL(inputURL string) bool {
	u, err := url.Parse(inputURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// doGetTokenRequest gets the list of valid locations
func (lc *LatencyChecker) doGetTokenRequest() (int, error) {

	if lc.GetAPIKey() == "NOT_SET" {
		return -1, errors.New(" DDOSIFY_X_API_KEY env var not set")
	}
	req, err := http.NewRequest(http.MethodGet, DDOSIFY_TOKEN_API_URL, nil)
	if err != nil {
		panic(err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req.Header.Add("X-API-KEY", lc.GetAPIKey())
	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return -1, err
	}

	bodyResponse := &tokenAPIResponse{}
	derr := json.NewDecoder(res.Body).Decode(bodyResponse)
	if derr != nil {
		return -2, derr
	}

	if res.StatusCode != http.StatusOK {
		return -3, errors.New(" Status code received: " + strconv.Itoa(res.StatusCode) + " ...but status code expected: " + strconv.Itoa(http.StatusOK))
	}

	return bodyResponse.RequestCount, nil

}

//doPostLatencyCheckRequest
func (lc *LatencyChecker) doPostLatencyCheckRequest() (map[string]interface{}, error) {
	var responseLatency map[string]interface{}
	reqBody := &latencyAPIRequest{TargetURL: lc.GetTargetURL(), Locations: lc.GetLocations()}
	reqBodyJson, _ := json.Marshal(reqBody)
	body := strings.NewReader(string(reqBodyJson))
	req, err := http.NewRequest(http.MethodPost, DDOSIFY_LATENCY_API_URL, body)
	if err != nil {
		panic(err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req.Header.Add("X-API-KEY", lc.GetAPIKey())
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	derr := json.NewDecoder(res.Body).Decode(&responseLatency)
	if derr != nil {
		return nil, derr
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("[ERROR] Status code received: " + strconv.Itoa(res.StatusCode) + " ...but status code expected: " + strconv.Itoa(http.StatusOK))
	}
	return responseLatency, nil
}
