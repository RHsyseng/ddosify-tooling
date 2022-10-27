/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/RHsyseng/ddosify-tooling/tooling/pkg/ddosify"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeReady                 string = "Ready"
	ConditionTypeLatencyCheckerRunning string = "LatencyCheckerRunning"
	ConditionIntervalTimeValid         string = "IntervalTimeValid"
	ConditionAPITokenValid             string = "APITokenValid"
)

type LatencyCheckerProvider struct {
	ProviderName string `json:"providerName"`
	APIKey       string `json:"apiKey"`
}

// LatencyCheckSpec defines the desired state of LatencyCheck
type LatencyCheckSpec struct {
	TargetURL             string                 `json:"targetURL"`
	NumberOfRuns          int                    `json:"numberOfRuns"`
	WaitInterval          string                 `json:"waitInterval"`
	Locations             []string               `json:"locations"`
	OutputLocationsNumber int                    `json:"outputLocationsNumber"`
	Provider              LatencyCheckerProvider `json:"provider"`
	Scheduled             bool                   `json:"scheduled,omitempty"`
	ScheduleDefinition    string                 `json:"scheduleDefinition,omitempty"`
}

type LatencyCheckResult struct {
	ExecutionTime string                            `json:"executionTime"`
	Result        *ddosify.LatencyCheckerOutputList `json:"result,omitempty"`
}

// LatencyCheckStatus defines the observed state of LatencyCheck
type LatencyCheckStatus struct {
	Results       LatencyCheckResult `json:"results"`
	LastExecution string             `json:"lastExecution"`
	Conditions    []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LatencyCheck is the Schema for the latencychecks API
type LatencyCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LatencyCheckSpec   `json:"spec,omitempty"`
	Status LatencyCheckStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LatencyCheckList contains a list of LatencyCheck
type LatencyCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LatencyCheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LatencyCheck{}, &LatencyCheckList{})
}
