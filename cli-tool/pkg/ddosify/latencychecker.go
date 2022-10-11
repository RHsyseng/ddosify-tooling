package ddosify

//LatencyChecker is the basic object struct to define the checker attributes
type LatencyChecker struct {
	TargetUrl string //This is the basic Target URL to be used
}

//NewLatencyChecker is the object constructor
func NewLatencyChecker(targetURL string) *LatencyChecker {
	return &LatencyChecker{TargetUrl: targetURL}
}

//GetTargetURL is the get method to retrieve the TargetURL attribute
func (lc LatencyChecker) GetTargetURL() string {
	return lc.TargetUrl
}
