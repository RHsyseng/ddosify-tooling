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
	log.Printf("Required tokens for this exection %d, available tokens: %d", requiredTokens, availableTokens)
	if availableTokens < requiredTokens {
		return LatencyCheckerOutputList{}, errors.New(" insufficient tokens")
	}
	// We need to wait before sending a new API request, otherwise we get throttled
	time.Sleep(DDOSIFY_API_THROTTLER_TIME * time.Second)

	latencyResults := make(map[string]float64)

	for i := 1; i <= lc.GetRuns(); i++ {
		// TODO: add progress bar
		// TODO: wait the interval time
		// TODO: validate the interval time parameter (investigate if we can delegate the validation to cobra somehow)
		log.Printf("Run number %d", i)
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
