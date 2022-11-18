package ddosify

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/adhocore/gronx"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

//GetEnv returns the value for a given Env Var
func GetEnv(varName string, defaultValue string) string {
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

func ValidateIntervalTime(interval string) bool {
	r := regexp.MustCompile(`^(\d*)(s|m|h)`)
	matched := r.MatchString(interval)
	return matched
}

func ValidateCronTime(cronTime string) bool {
	log.Println("[INFO] crontime to be validated: ", cronTime)
	gron := gronx.New()
	if cronTime != "" {
		return gron.IsValid(cronTime)
	}
	return false
}

func GetNextTimeCronTime(cronTime string) int64 {
	// We get the time of the next execution based on the scheduled
	nextTime, err := gronx.NextTick(cronTime, false)
	if err != nil {
		log.Println("[INFO] CronTime format is not valid or not scheduled cr passed: ", err)
		return -1
	}
	// We convert times to epoch and get the seconds remaining to the nextTime
	duration := nextTime.Unix() - time.Now().Unix()
	return duration
}

func IntervalTimeToSeconds(interval string) int {
	r := regexp.MustCompile(`^(\d*)(s|m|h)`)
	captureGroups := r.FindStringSubmatch(interval)
	if len(captureGroups) < 1 {
		return -1
	}
	timeValue, err := strconv.Atoi(captureGroups[1])
	if err != nil {
		return -1
	}
	timeUnit := captureGroups[2]

	switch timeUnit {
	case "s":
		return timeValue
	case "m":
		return timeValue * 60
	case "h":
		return timeValue * 3600
	default:
		return 1
	}
}

// doGetTokenRequest gets the list of valid locations
func (lc *LatencyChecker) doGetTokenRequest() (int, error) {

	if lc.GetAPIKey() == "NOT_SET" {
		return -1, errors.New(" DDOSIFY_X_API_KEY env var not set")
	}
	req, err := http.NewRequest(http.MethodGet, lc.GetServiceAPITokenURL(), nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-API-KEY", lc.GetAPIKey())
	res, _ := http.DefaultClient.Do(req)

	bodyResponse := &tokenAPIResponse{}
	derr := json.NewDecoder(res.Body).Decode(bodyResponse)
	if derr != nil {
		return -2, derr
	}

	if res.StatusCode != http.StatusOK {
		return -3, errors.New(" Status code received: " + strconv.Itoa(res.StatusCode) + " ...but status code expected: " + strconv.Itoa(http.StatusOK))
	}
	defer res.Body.Close()
	return bodyResponse.RequestCount, nil

}

//doPostLatencyCheckRequest runs a latency check and returns its result
func (lc *LatencyChecker) doPostLatencyCheckRequest() (map[string]interface{}, error) {

	var responseLatency map[string]interface{}
	reqBody := LatencyAPIRequest{
		TargetURL: lc.TargetUrl,
		Locations: lc.GetLocations(),
	}

	reqBodyJson, _ := json.Marshal(reqBody)

	body := bytes.NewReader(reqBodyJson)

	req, err := http.NewRequest(http.MethodPost, lc.GetServiceAPIURL(), body)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", CONTENT_TYPE_REQ)
	req.Header.Add("X-API-KEY", lc.GetAPIKey())

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	derr := json.NewDecoder(res.Body).Decode(&responseLatency)
	if derr != nil {
		log.Println(derr.Error())
		return nil, derr
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(" Status code received: " + strconv.Itoa(res.StatusCode) + " ...but status code expected: " + strconv.Itoa(http.StatusOK))
	}

	return responseLatency, nil
}

//getMinimumLatencies returns the minimum avg latencies from a latencychecktest result
func (lc *LatencyChecker) getMinimumLatencies(latencies map[string]float64) ([]string, []float64) {
	outputKeys := make([]string, len(latencies))
	outputLatency := make([]float64, len(latencies))
	keys := make([]string, 0, len(latencies))
	for k := range latencies {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return latencies[keys[i]] < latencies[keys[j]]
	})
	// If more output than available latencies, we fix to the number of latencies
	if lc.GetOutputLocationsNumber() > len(latencies) {
		lc.SetOutputLocationsNumber(len(latencies))
	}

	for i := 0; i < lc.GetOutputLocationsNumber(); i++ {
		outputKeys[i] = keys[i]
		outputLatency[i] = latencies[keys[i]]
	}
	return outputKeys, outputLatency
}
