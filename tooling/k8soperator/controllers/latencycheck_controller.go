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

package controllers

import (
	"context"
	latencyv1alpha1 "github.com/RHsyseng/ddosify-tooling/tooling/k8soperator/api/v1alpha1"
	"github.com/RHsyseng/ddosify-tooling/tooling/pkg/ddosify"
	"github.com/go-logr/logr"
	acmPRV1 "github.com/open-cluster-management/multicloud-operators-placementrule/pkg/apis/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"
)

// LatencyCheckReconciler reconciles a LatencyCheck object
type LatencyCheckReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Finalizer for our objects
const (
	latencyCheckerFinalizer = "finalizer.latency.redhat.com"
	concurrentReconciles    = 10
)

//+kubebuilder:rbac:groups=latency.redhat.com,resources=latencychecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=latency.redhat.com,resources=latencychecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=latency.redhat.com,resources=latencychecks/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps.open-cluster-management.io,resources=placementrules,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

func (r *LatencyCheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconcile loop started")
	// Fetch the Latency instance
	instance := &latencyv1alpha1.LatencyCheck{}

	// This uses the API
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("LatencyCheck resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get LatencyCheck resources")
		return ctrl.Result{}, err
	}

	// Check if the CR is marked to be deleted
	isInstanceMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isInstanceMarkedToBeDeleted {
		log.Info("Instance marked for deletion, running finalizers")
		if contains(instance.GetFinalizers(), latencyCheckerFinalizer) {
			// Run the finalizer logic
			err := r.finalizeLatencyCheck(log, instance)
			if err != nil {
				// Don't remove the finalizer if we failed to finalize the object
				return ctrl.Result{}, err
			}
			log.Info("Instance finalizers completed")
			// Remove finalizer once the finalizer logic has run
			controllerutil.RemoveFinalizer(instance, latencyCheckerFinalizer)
			err = r.Update(ctx, instance)
			if err != nil {
				// If the object update fails, requeue
				return ctrl.Result{}, err
			}
		}
		log.Info("Instance can be deleted now")
		return ctrl.Result{}, nil
	}

	// Add Finalizers to the CR
	if !contains(instance.GetFinalizers(), latencyCheckerFinalizer) {
		if err := r.addFinalizer(log, instance, ctx); err != nil {
			log.Info("Error adding finalizer")
			return ctrl.Result{}, err
		}
		// We need to requeue after adding the finalizer. Keep in mind that generation shouldn't change so no reconcile will happen next
		// this may not be true if the API is mutating the object after creation (example omitempty boolean fields set to false will be deleted by the API).
		return ctrl.Result{Requeue: true}, nil
	}

	// Run LatencyChecks

	// This instance is an infinite run
	if instance.Spec.Scheduled {
		// long run
		output, err := r.runLatencyChecker(log, instance)
		log.Info("Long-lived run")
		// If error, the status will be empty, and we will requeue in case next time it goes well
		r.prepareLatencyCheckerStatus(log, err, instance, &output)
		// We need to update the status before the ACM integration since that relies on the status of the object
		r.updateLatencyCheckStatus(instance, log)
		// Check if there is an RHACM Integration and create the placementrule
		if !reflect.DeepEqual(instance.Spec.ACMIntegration, latencyv1alpha1.LatencyCheckerACMIntegration{}) && err == nil {
			log.Info("Detected ACM Integration")
			err = r.generateACMIntegration(log, instance)
			// Update status to set ACM conditions
			r.updateLatencyCheckStatus(instance, log)
			if err != nil {
				log.Info(err.Error())
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{RequeueAfter: time.Duration(ddosify.GetNextTimeCronTime(instance.Spec.ScheduleDefinition)) * time.Second, Requeue: true}, nil
	}

	log.Info("Short-lived run")
	if instance.Status.Results != nil {
		// If a short-lived run has a status it means it was executed already
		// we skip the reconcile, this can happen when the controller is restarted
		// or the spec is modified afterwards
		log.Info("Short-lived run already executed, skipping")
		return ctrl.Result{}, nil
	}
	output, err := r.runLatencyChecker(log, instance)

	// If error, the status will be empty, and we will requeue in case next time it goes well
	r.prepareLatencyCheckerStatus(log, err, instance, &output)
	// We need to update the status before the ACM integration since that relies on the status of the object
	r.updateLatencyCheckStatus(instance, log)
	// If error, the status will be empty
	if !reflect.DeepEqual(instance.Spec.ACMIntegration, latencyv1alpha1.LatencyCheckerACMIntegration{}) && err == nil {
		log.Info("Detected ACM Integration")
		err = r.generateACMIntegration(log, instance)
		// Update status to set ACM conditions
		r.updateLatencyCheckStatus(instance, log)
		if err != nil {
			log.Info(err.Error())
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil

}

func (r *LatencyCheckReconciler) prepareLatencyCheckerStatus(log logr.Logger, errRun error, instance *latencyv1alpha1.LatencyCheck, result *ddosify.LatencyCheckerOutputList) {

	newResult := latencyv1alpha1.LatencyCheckResult{
		ExecutionTime: time.Now().Format(time.RFC3339),
		Result:        result,
	}
	// We need to concatenate existing results to the new result
	instance.Status.Results = append(instance.Status.Results, newResult)
	instance.Status.LastExecution = time.Now().Format(time.RFC3339)
	if instance.Spec.Scheduled {
		instance.Status.NextExecution = time.Now().Add(time.Duration(ddosify.GetNextTimeCronTime(instance.Spec.ScheduleDefinition)) * time.Second).Format(time.RFC3339)
	}
	if errRun != nil {
		log.Info("Error running LatencyChecker")
		// Update status
		instance.Status.Results[0] = latencyv1alpha1.LatencyCheckResult{
			ExecutionTime: time.Now().Format(time.RFC3339),
			Result: &ddosify.LatencyCheckerOutputList{
				Result: []ddosify.LatencyCheckerOutput{},
			},
		}
		instance.Status.LastExecution = time.Now().Format(time.RFC3339)

		switch {
		case errors.IsBadRequest(errRun) && errRun.Error() == latencyv1alpha1.ConditionScheduleDefinitionValid:
			meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionScheduleDefinitionValid, Status: metav1.ConditionFalse, Reason: latencyv1alpha1.ConditionScheduleDefinitionValid, Message: latencyv1alpha1.ConditionScheduleDefinitionNotValidMsg})
			break
		case errors.IsBadRequest(errRun) && errRun.Error() == latencyv1alpha1.ConditionIntervalTimeValid:
			meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionIntervalTimeValid, Status: metav1.ConditionFalse, Reason: latencyv1alpha1.ConditionIntervalTimeValid, Message: latencyv1alpha1.ConditionIntervalTimeNotValidMsg})
			break
		case errors.IsInternalError(errRun):
			meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionAPITokenValid, Status: metav1.ConditionFalse, Reason: latencyv1alpha1.ConditionAPITokenValid, Message: latencyv1alpha1.ConditionAPITokenNotValidMsg})
			break
		}
		//set conditions
		meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionNotReady, Status: metav1.ConditionTrue, Reason: latencyv1alpha1.ConditionNotReady, Message: latencyv1alpha1.ConditionNotReadyMsg})
		return
	}
	meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionReady, Status: metav1.ConditionTrue, Reason: latencyv1alpha1.ConditionReady, Message: latencyv1alpha1.ConditionReadyMsg})
}

func (r *LatencyCheckReconciler) runLatencyChecker(log logr.Logger, cr *latencyv1alpha1.LatencyCheck) (ddosify.LatencyCheckerOutputList, error) {
	log.Info("About to run latency check")
	if !ddosify.ValidateIntervalTime(cr.Spec.WaitInterval) {
		log.Info("Invalid wait interval")
		return ddosify.LatencyCheckerOutputList{}, errors.NewBadRequest(latencyv1alpha1.ConditionIntervalTimeValid)
	}

	if cr.Spec.Scheduled && !ddosify.ValidateCronTime(cr.Spec.ScheduleDefinition) {
		log.Info("Invalid cron time")
		return ddosify.LatencyCheckerOutputList{}, errors.NewBadRequest(latencyv1alpha1.ConditionScheduleDefinitionValid)
	}

	lc := ddosify.NewLatencyChecker(cr.Spec.Provider.APIKey, cr.Spec.TargetURL, cr.Spec.NumberOfRuns, ddosify.IntervalTimeToSeconds(cr.Spec.WaitInterval), cr.Spec.Locations, cr.Spec.OutputLocationsNumber)
	res, err := lc.RunCommandExec()
	if err != nil {

		return ddosify.LatencyCheckerOutputList{}, errors.NewInternalError(err)
	}
	return res, nil
}

// addFinalizer adds a given finalizer to a given CR
func (r *LatencyCheckReconciler) addFinalizer(log logr.Logger, cr *latencyv1alpha1.LatencyCheck, ctx context.Context) error {
	log.Info("Adding Finalizer for the LatencyCheck")
	controllerutil.AddFinalizer(cr, latencyCheckerFinalizer)

	// Update CR
	err := r.Update(ctx, cr)
	if err != nil {
		log.Error(err, "Failed to update LatencyCheck with finalizer")
		return err
	}
	return nil
}

// updateLatencyCheckStatus updates the Status of a given CR
func (r *LatencyCheckReconciler) updateLatencyCheckStatus(cr *latencyv1alpha1.LatencyCheck, log logr.Logger) (*latencyv1alpha1.LatencyCheck, error) {
	latencyCheck := &latencyv1alpha1.LatencyCheck{}
	err := r.Get(context.Background(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, latencyCheck)
	if err != nil {
		return latencyCheck, err
	}

	if !reflect.DeepEqual(cr.Status, latencyCheck.Status) {
		log.Info("Updating LatencyCheck Status.")
		// We need to update the status
		err = r.Status().Update(context.Background(), cr)
		if err != nil {
			log.Info(err.Error())
			return cr, err
		}
		updatedlatencyCheck := &latencyv1alpha1.LatencyCheck{}
		err = r.Get(context.Background(), types.NamespacedName{Name: cr.Name, Namespace: cr.Namespace}, updatedlatencyCheck)
		if err != nil {
			return cr, err
		}
		cr = updatedlatencyCheck.DeepCopy()
	}
	return cr, nil

}

// finalizeLatencyCheck runs required tasks before deleting the objects owned by the CR
func (r *LatencyCheckReconciler) finalizeLatencyCheck(log logr.Logger, cr *latencyv1alpha1.LatencyCheck) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.
	log.Info("Successfully finalized LatencyCheck")
	return nil
}
func (r *LatencyCheckReconciler) generateACMIntegration(log logr.Logger, cr *latencyv1alpha1.LatencyCheck) error {
	placementRule := r.newPlacementRule(log, cr)
	// Check if placementrule already exists
	placementRuleFound := &acmPRV1.PlacementRule{}
	err := r.Get(context.Background(), types.NamespacedName{Name: placementRule.Name, Namespace: placementRule.Namespace}, placementRuleFound)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new PlacementRule", "PlacementRule.Namespace", placementRule.Namespace, "PlacementRule.Name", placementRule.Name)
		err = r.Create(context.Background(), placementRule)
		if err != nil {
			return err
		}
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionACMPlacementRuleCreated, Status: metav1.ConditionTrue, Reason: latencyv1alpha1.ConditionACMPlacementRuleCreated, Message: latencyv1alpha1.ConditionACMPlacementRuleCreatedMsg})
	} else if err != nil {
		return err
	} else {
		log.Info("PlacementRule updated", "placementRule.Namespace", placementRuleFound.Namespace, "placementRule.Name", placementRuleFound.Name)
		placementRuleFound.Spec.ClusterReplicas = &cr.Spec.ACMIntegration.PlacementRuleClusterReplicas
		placementRuleFound.Spec.ClusterSelector.MatchLabels[cr.Spec.ACMIntegration.ClusterLocationLabel] = r.getPlacementRuleLocation(cr)
		err = r.Update(context.Background(), placementRuleFound)
		if err != nil {
			return err
		}
		meta.SetStatusCondition(&cr.Status.Conditions, metav1.Condition{Type: latencyv1alpha1.ConditionACMPlacementRuleUpdated, Status: metav1.ConditionTrue, Reason: latencyv1alpha1.ConditionACMPlacementRuleUpdated, Message: latencyv1alpha1.ConditionACMPlacementRuleUpdatedMsg})
	}
	return nil
}

func (r *LatencyCheckReconciler) getPlacementRuleLocation(cr *latencyv1alpha1.LatencyCheck) string {
	// len(cr.Status.Results)-1 will return latest execution result, Result[0].Location will return the best location
	latencyCheckBestLocation := cr.Status.Results[len(cr.Status.Results)-1].Result.Result[0].Location
	switch cr.Spec.ACMIntegration.LocationMatchingStrategy {
	case "continent":
		// NA
		return strings.Split(latencyCheckBestLocation, ".")[0]
	case "country":
		// NA.US
		return strings.Join(strings.Split(latencyCheckBestLocation, ".")[:2], ".")
	case "state":
		// NA.US.PA
		return strings.Join(strings.Split(latencyCheckBestLocation, ".")[:3], ".")
		// case "city" will be covered by default case
	default:
		// NA.US.PA.PH
		return latencyCheckBestLocation
	}
}

func (r *LatencyCheckReconciler) newPlacementRule(log logr.Logger, cr *latencyv1alpha1.LatencyCheck) *acmPRV1.PlacementRule {
	bestLocation := r.getPlacementRuleLocation(cr)
	labels := map[string]string{
		"app": cr.Name,
	}
	clusterSelectorLabels := map[string]string{

		cr.Spec.ACMIntegration.ClusterLocationLabel: bestLocation,
	}
	if cr.Spec.ACMIntegration.PlacementRuleNamespace == "" {
		cr.Spec.ACMIntegration.PlacementRuleNamespace = cr.Namespace
	} else {
		placementRuleNamespaceFound := &corev1.Namespace{}
		err := r.Get(context.Background(), types.NamespacedName{Name: cr.Spec.ACMIntegration.PlacementRuleNamespace, Namespace: ""}, placementRuleNamespaceFound)
		if err != nil && errors.IsNotFound(err) {
			log.Info("Namespace not found, defaulting to cr namespace")
			cr.Spec.ACMIntegration.PlacementRuleNamespace = cr.Namespace
		}
	}
	if cr.Spec.ACMIntegration.PlacementRuleName == "" {
		cr.Spec.ACMIntegration.PlacementRuleName = "placementrule-" + cr.Name
	}

	if cr.Spec.ACMIntegration.ClusterLocationLabel == "" {
		cr.Spec.ACMIntegration.ClusterLocationLabel = "region"
	}

	// If no replicas are defined in the spec, it will be set to 1

	return &acmPRV1.PlacementRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps.open-cluster-management.io/v1",
			Kind:       "PlacementRule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.ACMIntegration.PlacementRuleName,
			Namespace: cr.Spec.ACMIntegration.PlacementRuleNamespace,
			Labels:    labels,
		},
		Spec: acmPRV1.PlacementRuleSpec{
			ClusterReplicas: &cr.Spec.ACMIntegration.PlacementRuleClusterReplicas,
			GenericPlacementFields: acmPRV1.GenericPlacementFields{
				ClusterSelector: &metav1.LabelSelector{
					MatchLabels: clusterSelectorLabels,
				},
			},
			ClusterConditions: nil,
		},
		Status: acmPRV1.PlacementRuleStatus{
			Decisions: nil,
		},
	}
}

// contains returns true if a string is found on a slice
func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// Ignore changes that do not increase the resource generation
func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *LatencyCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&latencyv1alpha1.LatencyCheck{}).
		WithEventFilter(ignoreDeletionPredicate()).
		WithOptions(controller.Options{MaxConcurrentReconciles: concurrentReconciles}).
		Complete(r)
}
