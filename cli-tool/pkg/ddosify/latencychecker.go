package ddosify

//LatencyChecker is the basic object struct to define the checker attributes
const (
	CONTENT_TYPE_REQ            = "application/json"
	DDOSIFY_TOKEN_API_URL       = "https://api.ddosify.com/v1/balance"
	DDOSIFY_LATENCY_API_URL     = "https://api.ddosify.com/v1/latency/test"
	DDOSIFY_LATENCY_TOKENS_COST = 5000
	DDOSIFY_API_THROTTLER_TIME  = 1
)

type LatencyChecker struct {
	TargetUrl             string   //This is the basic Target URL to be used
	Runs                  int      //This is the number of executions to be run
	WaitInterval          int      //This is the amount of time to wait between runs (in seconds)
	Locations             []string //This is the locations to be used for the latency check
	APIKey                string   //This is the API Key to authenticate to ddosify (we get this from an ENV var)
	ContentType           string   //This is the content type used to interact with the ddosify API
	OutputLocationsNumber int      //This is the number of locations to be outputed
	ServiceAPITokenURL    string
	ServiceAPIURL         string
}

type LatencyCheckerOutput struct {
	Location   string  `json:"location",yaml:"location"`
	AvgLatency float64 `json:"avgLatency",yaml:"avgLatency"`
}

type LatencyCheckerOutputList struct {
	Result []LatencyCheckerOutput `json:"result",yaml:"result"`
}

type tokenAPIResponse struct {
	RequestCount int `json:"request_count"`
	Duration     int `json:"duration"`
}

type LatencyAPIRequest struct {
	TargetURL string   `json:"target"`
	Locations []string `json:"locations"`
}

//NewLatencyChecker is the object constructor
func NewLatencyChecker(targetURL string, runs int, waitInterval int, locations []string, outputLocationsNumber int) *LatencyChecker {
	return &LatencyChecker{
		TargetUrl:             targetURL,
		Runs:                  runs,
		WaitInterval:          waitInterval,
		Locations:             locations,
		APIKey:                getEnv("DDOSIFY_X_API_KEY", "NOT_SET"),
		ContentType:           CONTENT_TYPE_REQ,
		OutputLocationsNumber: outputLocationsNumber,
		ServiceAPITokenURL:    DDOSIFY_TOKEN_API_URL,
		ServiceAPIURL:         DDOSIFY_LATENCY_API_URL,
	}
}

//GetServiceAPITokenURLis the get method to retrieve the ServiceAPITokenURL attribute
func (lc *LatencyChecker) GetServiceAPITokenURL() string {
	return lc.ServiceAPITokenURL
} //GetServiceAPIURL is the get method to retrieve the ServiceAPIURL attribute
func (lc *LatencyChecker) GetServiceAPIURL() string {
	return lc.ServiceAPIURL
}

//GetTargetURL is the get method to retrieve the TargetURL attribute
func (lc *LatencyChecker) GetTargetURL() string {
	return lc.TargetUrl
}

//GetRuns is the get method to retrieve the Runs attribute
func (lc *LatencyChecker) GetRuns() int {
	return lc.Runs
}

//GetWaitInterval is the get method to retrieve the WaitInterval attribute
func (lc *LatencyChecker) GetWaitInterval() int {
	return lc.WaitInterval
}

//GetLocations is the get method to retrieve the WaitInterval attribute
func (lc *LatencyChecker) GetLocations() []string {
	return lc.Locations
}

//GetAPIKey is the method to retrieve the APIKey attribute
func (lc *LatencyChecker) GetAPIKey() string {
	return lc.APIKey
}

//GetOutputLocationsNumber is the get method to retrieve the OutputLocationsNumber attribute
func (lc *LatencyChecker) GetOutputLocationsNumber() int {
	return lc.OutputLocationsNumber
}

//SetTargetURL is the setter method to set the TargetURL attribute
func (lc *LatencyChecker) SetTargetURL(targetUrl string) {
	lc.TargetUrl = targetUrl
}

//SetRuns is the setter method to set the Runs attribute
func (lc *LatencyChecker) SetRuns(runs int) {
	lc.Runs = runs
}

//SetWaitInterval is the setter method to set the WaitInterval attribute
func (lc *LatencyChecker) SetWaitInterval(interval int) {
	lc.WaitInterval = interval
}

//SetLocations is the setter method to set the WaitInterval attribute
func (lc *LatencyChecker) SetLocations(locations []string) {
	lc.Locations = locations
}

//SetOutputLocationsNumber is the setter method to set the Runs attribute
func (lc *LatencyChecker) SetOutputLocationsNumber(outputLocationsNumber int) {
	lc.OutputLocationsNumber = outputLocationsNumber
}

//SetServiceAPITokenURL is the setter method to set the Runs attribute
func (lc *LatencyChecker) SetServiceAPITokenURL(URL string) {
	lc.ServiceAPITokenURL = URL
}

// SetServiceAPIURL is the setter method to set the Runs attribute
func (lc *LatencyChecker) SetServiceAPIURL(URL string) {
	lc.ServiceAPIURL = URL
}
