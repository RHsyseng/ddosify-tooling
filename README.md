# Ddosify latency tooling
![Coverage](https://img.shields.io/badge/Coverage-71.2%25-brightgreen)

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

## K8s Controller

Following the [cli](./tooling/cmd/) we have coded a [Kubernetes Controller](./tooling/k8soperator/), this operator has the same features exposed in the CLI, plus a scheduling mechanism.

This controller implements the `latencycheck.latency.redhat.com` API, this API is the user interface to schedule/run latencychecks using Kubernetes.

Users can leverage this API to run one-shot latency checks or to schedule latency checks over time.

Below a one-shot example:

:information_source: Below configuration defines a latencycheck that runs only once.

~~~yaml
apiVersion: latency.redhat.com/v1alpha1
kind: LatencyCheck
metadata:
  name: lc-shortlived
spec:
  targetURL: "https://google.com"
  numberOfRuns: 1
  waitInterval: "5s"
  locations:
  - "NA.*"
  outputLocationsNumber: 3
  scheduled: false
  provider:
    providerName: "ddosify"
    apiKey: "90738bc2-9073-4c93-8bc2-ae7608df1b1e"
~~~

We can also define a scheduled latencycheck:

:information_source: Below configuration defines a latencycheck that runs every minute. The `scheduleDefinition` parameter supports cron-like scheduling definitions.

~~~yaml
apiVersion: latency.redhat.com/v1alpha1
kind: LatencyCheck
metadata:
  name: lc-longlived
spec:
  targetURL: "https://google.com"
  numberOfRuns: 2
  waitInterval: "10s"
  locations:
  - "NA.*"
  outputLocationsNumber: 3
  scheduled: true
  scheduleDefinition: "*/1 * * * *"
  provider:
    providerName: "ddosify"
    apiKey: "90738bc2-9073-4c93-8bc2-ae7608df1b1e"
~~~

The way the controller reports the results to the user is via the `status` field, you can see an example below:

~~~yaml
status:
  conditions:
  - lastTransitionTime: "2022-11-18T15:14:53Z"
    message: waitInterval is not valid
    reason: IntervalTimeValid
    status: "False"
    type: IntervalTimeValid
  lastExecution: "2022-11-18T16:14:53+01:00"
  results:
  - execution:
      result:
      - avgLatency: 47
        location: NA.US.WA.SE
      - avgLatency: 48
        location: NA.US.WA.QU
      - avgLatency: 48
        location: NA.US.SC.NC
    executionTime: "2022-11-18T16:14:53+01:00"
~~~

There are still some rough edges that will be addresses in the next release.

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