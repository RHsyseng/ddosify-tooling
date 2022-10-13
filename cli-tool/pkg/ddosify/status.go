package ddosify

import (
	"log"
)

func (lc *LatencyChecker) RunCommandStatus() error {

	log.Printf("TargetURL %s, Number of Runs: %d, Wait interval: %s", lc.GetTargetURL(), lc.GetRuns(), lc.GetWaitInterval())
	return nil
}

func (lc *LatencyChecker) execStatus() error {

	return nil
}
