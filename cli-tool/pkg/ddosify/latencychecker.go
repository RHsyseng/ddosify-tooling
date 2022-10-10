package ddosify

type LatencyChecker struct {
	TargetUrl string
}

func NewLatencyChecker(targetURL string) *LatencyChecker {
	return &LatencyChecker{TargetUrl: targetURL}
}

func (lc LatencyChecker) GetTargetURL() string {
	return lc.TargetUrl
}
