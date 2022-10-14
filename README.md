# Ddosify latency tooling

Currently, ddosify exposes an API that can be used to test latencies, the documentation can be found [here](https://docs.ddosify.com/cloud/api/latency-testing-api). It has some limitations:

* Cannot do scheduled runs
* There is no client for ddosify cloud

## CLI Tool

The CLI tool interacts with the [ddosify latency API](https://docs.ddosify.com/cloud/api/latency-testing-api). The CLI usage can be found below:

~~~sh
Usage:
  ddosify-latencies run [flags]

Flags:
  -h, --help                    help for run
  -i, --interval string         The amount of waiting time between runs. (default "1m")
  -l, --locations stringArray   The array of locations to be requested. e.g: NA.US.*,NA.EU.* (default [EU.ES.*])
  -o, --output-format string    Output in an specific format. Usage: '-o [ table | yaml | json ]' (default "table")
      --output-locations int    The number of best locations to output. (default 1)
  -r, --runs int                The number of executions. (default 1)
  -t, --target-url string       The target url. e.g: https://google.com
~~~

An example output using `table` output:

~~~sh
┌─────────────┬─────────────────┐
│ LOCATION    │ AVERAGE LATENCY │
├─────────────┼─────────────────┤
│ NA.US.TX.DA │ 9.000000        │
│ NA.US.TX.SA │ 11.000000       │
│ NA.US.NV.LV │ 13.000000       │
└─────────────┴─────────────────┘
~~~

An example output using `json` output:

~~~sh
{
    "result": [
        {
            "location": "NA.US.TX.HO",
            "avgLatency": 3
        },
        {
            "location": "NA.US.VA.AS",
            "avgLatency": 4
        },
        {
            "location": "NA.US.TX.DA",
            "avgLatency": 4
        }
    ]
}
~~~

An example output using `yaml` output:

~~~sh
result:
- location: NA.US.VA.AS
  avglatency: 4
- location: NA.US.TX.DA
  avglatency: 4
- location: NA.US.TX.HO
  avglatency: 4
~~~

## Idea 2 - K8s Controller

This is a bit more complicated to code, but still doable in a short time. We need to define a new CRD, for example, `latencychecks.ddosify.com`. This CRD will have the following spec fields:

* Target URL - String
* Locations (as in [docs](https://docs.ddosify.com/cloud/api/latency-testing-api#example-usages-of-locations-object)) - Slice
* Number of runs - Int
* Time between runs - String? (to support 60s, 1m, 1d, 1M, 1y…)
* Recurrent - Bool

For the CRD status we will have:

* Best location
* Last run time
* Last run status
* Next run time?

The idea is a controller that runs the checks and outputs the best location to its CRD status.

## Idea 3 - K8s Controller Integration with RHACM

This will be the most complicated idea to code, but probably the one that provides a more compelling user story.

We will use the output from idea2 to update/create a PlacementRule.

Prerequisites:

* Managed clusters need to have labels matching the locations (as in [docs](https://docs.ddosify.com/cloud/api/latency-testing-api#example-usages-of-locations-object))

New spec fields:

* PlacementRuleNamespace
* PlacementRuleName

New status fields:

* PlacementRuleAction - String (Created / Updated / NoAction)