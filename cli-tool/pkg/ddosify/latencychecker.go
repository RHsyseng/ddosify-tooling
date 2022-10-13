package ddosify

//LatencyChecker is the basic object struct to define the checker attributes
const (
	CONTENT_TYPE_REQ            = "application/json"
	DDOSIFY_TOKEN_API_URL       = "https://api.ddosify.com/v1/balance/"
	DDOSIFY_LATENCY_API_URL     = "https://api.ddosify.com/v1/latency/test/"
	DDOSIFY_LATENCY_TOKENS_COST = 5000
)

type LatencyChecker struct {
	TargetUrl    string   //This is the basic Target URL to be used
	Runs         int      //This is the number of executions to be run
	WaitInterval string   //This is the amount of time to wait between runs
	Locations    []string //This is the locations to be used for the latency check
	APIKey       string   //This is the API Key to authenticate to ddosify (we get this from an ENV var)
	ContentType  string   //This is the content type used to interact with the ddosify API
}

type tokenAPIResponse struct {
	RequestCount int `json:"request_count"`
	Duration     int `json:"duration"`
}

type latencyAPIRequest struct {
	TargetURL string   `json:"target"`
	Locations []string `json:"locations"`
}

//NewLatencyChecker is the object constructor
func NewLatencyChecker(targetURL string, runs int, waitInterval string, locations []string) *LatencyChecker {
	return &LatencyChecker{
		TargetUrl:    targetURL,
		Runs:         runs,
		WaitInterval: waitInterval,
		Locations:    locations,
		APIKey:       getEnv("DDOSIFY_X_API_KEY", "NOT_SET"),
		ContentType:  CONTENT_TYPE_REQ,
	}
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
func (lc *LatencyChecker) GetWaitInterval() string {
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

//SetTargetURL is the get method to retrieve the TargetURL attribute
func (lc *LatencyChecker) SetTargetURL(targetUrl string) {
	lc.TargetUrl = targetUrl
}

//SetRuns is the get method to retrieve the Runs attribute
func (lc *LatencyChecker) SetRuns(runs int) {
	lc.Runs = runs
}

//SetWaitInterval is the get method to retrieve the WaitInterval attribute
func (lc *LatencyChecker) SetWaitInterval(interval string) {
	lc.WaitInterval = interval
}

//SetLocations is the get method to retrieve the WaitInterval attribute
func (lc *LatencyChecker) SetLocations(locations []string) {
	lc.Locations = locations
}

// TODO: get tokens available https://api.ddosify.com/v1/balance/

// TODO: run latency check https://api.ddosify.com/v1/latency/test/
