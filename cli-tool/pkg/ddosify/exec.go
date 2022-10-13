package ddosify

import (
	"errors"
	"log"
)

func (lc *LatencyChecker) RunCommandExec() error {

	// Get number of runs, each run costs 5k tokens. Validate we have enough tokens to do all requests
	availableTokens, err := lc.doGetTokenRequest()
	if err != nil {
		switch availableTokens {
		case -1:
			log.Println("error detected when running the request to the Token API")
			break
		case -2:
			log.Println("error detected when trying to decode API response")
			break
		case -3:
			log.Println("unexpected http response code")
			break
		default:
			// This shouldn't happen
			log.Println("unexpected error")
			break
		}
		return err
	}

	// We have the available tokens
	requiredTokens := lc.GetRuns() * DDOSIFY_LATENCY_TOKENS_COST
	log.Printf("Required tokens for this exection %d, available tokens: %d", requiredTokens, availableTokens)
	if availableTokens < requiredTokens {
		return errors.New(" insufficient tokens")
	}

	for i := 1; i <= lc.GetRuns(); i++ {
		// TODO: add progress bar
		log.Printf("Run number %d", i)
		// Run the latency check
		responseLatencyCheck, err := lc.doPostLatencyCheckRequest()
		if err != nil {
			log.Println("Error doing Latency Check Request", err.Error())
		}
		for key, val := range responseLatencyCheck {
			log.Println(key)
			log.Println(val.(map[string]interface{})["avg_duration"])
			log.Println(val.(map[string]interface{})["status_code"])
		}

	}

	log.Printf("TargetURL %s, Number of Runs: %d, Wait interval: %s, Locations: %s", lc.GetTargetURL(), lc.GetRuns(), lc.GetWaitInterval(), lc.GetLocations())
	return nil
}

func (lc *LatencyChecker) execStatus() error {

	return nil
}
