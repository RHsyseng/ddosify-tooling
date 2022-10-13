package ddosify

//LatencyChecker is the basic object struct to define the checker attributes
type LatencyChecker struct {
	TargetUrl    string //This is the basic Target URL to be used
	Runs         int    //This is the number of executions to be run
	WaitInterval string //This is the amount of time to wait between runs
}

//NewLatencyChecker is the object constructor
func NewLatencyChecker(targetURL string, runs int, waitInterval string) *LatencyChecker {
	return &LatencyChecker{TargetUrl: targetURL, Runs: runs, WaitInterval: waitInterval}
}

//GetTargetURL is the get method to retrieve the TargetURL attribute
func (lc LatencyChecker) GetTargetURL() string {
	return lc.TargetUrl
}

//GetRuns is the get method to retrieve the Runs attribute
func (lc LatencyChecker) GetRuns() int {
	return lc.Runs
}

//GetWaitInterval is the get method to retrieve the WaitInterval attribute
func (lc LatencyChecker) GetWaitInterval() string {
	return lc.WaitInterval
}
