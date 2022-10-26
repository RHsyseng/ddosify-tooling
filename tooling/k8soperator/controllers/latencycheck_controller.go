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
	latencyv1alpha1 "github.com/RHsyseng/ddosify-tooling/api/v1alpha1"
	"github.com/RHsyseng/ddosify-tooling/tooling/pkg/ddosify"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
	//latency "github.com/RHsyseng/ddosify-tooling/tooling/pkg/ddosify"
)

// LatencyCheckReconciler reconciles a LatencyCheck object
type LatencyCheckReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Finalizer for our objects
const latencyCheckerFinalizer = "finalizer.latency.redhat.com"

//+kubebuilder:rbac:groups=latency.redhat.com,resources=latencychecks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=latency.redhat.com,resources=latencychecks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=latency.redhat.com,resources=latencychecks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LatencyCheck object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
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
		if err := r.addFinalizer(log, instance); err != nil {
			return ctrl.Result{}, err
		}
	}
	// Run LatencyChecks
	// Create LatencyCheckerObject

	log.Info("Add or Update")
	// This instance is an infinite run
	if instance.Spec.Scheduled {
		r.runLatencyChecker(log, instance)
		log.Info("Long-lived run")
		return ctrl.Result{RequeueAfter: 60 * time.Second, Requeue: true}, nil
	}
	r.runLatencyChecker(log, instance)
	log.Info("Short-lived run")
	return ctrl.Result{}, nil
}

func (r *LatencyCheckReconciler) runLatencyChecker(log logr.Logger, cr *latencyv1alpha1.LatencyCheck) (ddosify.LatencyCheckerOutputList, error) {
	log.Info("About to run latency check")

	lc := ddosify.NewLatencyChecker(cr.Spec.TargetURL, cr.Spec.NumberOfRuns, cr.Spec.WaitInterval, cr.Spec.Locations, cr.Spec.OutputLocationsNumber)
	res, err := lc.RunCommandExec()
	if err != nil {
		return ddosify.LatencyCheckerOutputList{}, err
	}
	return res, nil
}

// addFinalizer adds a given finalizer to a given CR
func (r *LatencyCheckReconciler) addFinalizer(log logr.Logger, cr *latencyv1alpha1.LatencyCheck) error {
	log.Info("Adding Finalizer for the LatencyCheck")
	controllerutil.AddFinalizer(cr, latencyCheckerFinalizer)

	// Update CR
	err := r.Update(context.Background(), cr)
	if err != nil {
		log.Error(err, "Failed to update LatencyCheck with finalizer")
		return err
	}
	return nil
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

// contains returns true if a string is found on a slice
func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *LatencyCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&latencyv1alpha1.LatencyCheck{}).
		Complete(r)
}
