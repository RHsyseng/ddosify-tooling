package ddosify

import (
	"log"
)

func RunCommandStatus() error {
	lc := NewLatencyChecker("")

	targetURL := lc.GetTargetURL()
	log.Printf("TargetURL %s", targetURL)
	return nil
}

func (lc *LatencyChecker) execStatus() error {

	return nil
}
