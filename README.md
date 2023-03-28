## CLI Tool

CLI has been migrated to https://github.com/ddosify/ddosify-latency-cli

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

## K8s Controller Integration with RHACM

The controller can be integrated with RHACM, the way it integrates with RHACM is via the creation of PlacementRules.

There is a new section in the `LatencyCheck` spec, this section is named `acmIntegration` and it looks like this:

~~~yaml
  acmIntegration:
    clusterLocationLabel: ddosify
    locationMatchingStrategy: city
    placementRuleClusterReplicas: 1
    placementRuleName: ddosify-demo
    placementRuleNamespace: DEMO_NAMESPACE
~~~

* `clusterLocationLabel` defines the label that is used in the managed clusters to match them to a ddosify location. It must use [ddosify locations](https://docs.ddosify.com/cloud/api/latency-testing-api).
* `locationMatchingStrategy` defines at what level we match the managed clusters, we can chose between continent, country, state or city. For example, for `NA.US.PA.PH`: Continent will be NA, country will be NA.US, state will be NA.US.PA, and city will be NA.US.PA.PH.
* `placementRuleClusterReplicas` defines the number of clusters we want to target in case we have more than one matching a specific location.
* `placementRuleName` defines the name for the PlacementRule that will be created or updated in case it already exists.
* `placementRuleNamespace` defines the namespace where the PlacementRule will be created.
