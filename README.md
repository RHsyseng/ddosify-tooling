# Ddosify latency tooling

Currently, ddosify exposes an API that can be used to test latencies, the documentation can be found [here](https://docs.ddosify.com/cloud/api/latency-testing-api). It has some limitations:

* Cannot do scheduled runs
* There is no client for ddosify cloud

## Idea 1 - CLI Tool

This is the most simple and easy to code idea, it will be a simple CLI tool where you can define:

* Target URL
* Locations (as in [docs](https://docs.ddosify.com/cloud/api/latency-testing-api#example-usages-of-locations-object))
* Number of runs
* Time between runs

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