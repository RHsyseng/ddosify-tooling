package ddosify

import (
	"errors"
	"log"
	"time"
)

func (lc *LatencyChecker) RunCommandExec() (LatencyCheckerOutputList, error) {

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
		return LatencyCheckerOutputList{}, err
	}

	// We have the available tokens
	requiredTokens := lc.GetRuns() * DDOSIFY_LATENCY_TOKENS_COST
	log.Printf("Required tokens for this execution %d, available tokens: %d", requiredTokens, availableTokens)
	if availableTokens < requiredTokens {
		return LatencyCheckerOutputList{}, errors.New(" insufficient tokens")
	}
	// We need to wait before sending a new API request, otherwise we get throttled
	time.Sleep(DDOSIFY_API_THROTTLER_TIME * time.Second)

	latencyResults := make(map[string]float64)

	log.Printf("Sleeping %ds between latency requests", lc.GetWaitInterval())

	for i := 1; i <= lc.GetRuns(); i++ {
		log.Printf("Request number [%d/%d]", i, lc.GetRuns())
		// Run the latency check
		responseLatencyCheck, err := lc.doPostLatencyCheckRequest()
		if err != nil {
			log.Println("Error doing Latency Check Request", err.Error())
		}

		for key, val := range responseLatencyCheck {
			latency := val.(map[string]interface{})["latency"]
			status_code := val.(map[string]interface{})["status_code"]

			// If a location fails, we want to penalize it
			if status_code.(float64) != 200 {
				latency = 1000
			}
			latencyResults[key] += latency.(float64)
		}
		if lc.GetRuns() > 1 {
			// Wait before running next iteration
			time.Sleep(time.Duration(lc.GetWaitInterval()) * time.Second)
		}

	}

	var outputList LatencyCheckerOutputList
	var output LatencyCheckerOutput

	bestLocation, avgLatencies := lc.getMinimumLatencies(latencyResults)
	for i := 0; i < lc.GetOutputLocationsNumber(); i++ {
		output.AvgLatency = avgLatencies[i] / float64(lc.GetRuns())
		output.Location = bestLocation[i]
		outputList.Result = append(outputList.Result, output)
	}
	return outputList, nil
}
