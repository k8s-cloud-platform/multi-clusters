/*
Copyright 2022 The MultiClusterPlatform Authors.

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

package scheduler

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "github.com/multi-cluster-platform/mcp/pkg/apis/apps/v1alpha1"
)

type Scheduler struct {
	client.Client
}

var _ reconcile.Reconciler = &Scheduler{}

// SetupWithManager sets up the controller with the Manager.
func (s *Scheduler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Deployable{}).
		Complete(s)
}

// Reconcile handles scheduleOne logic
func (s *Scheduler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	klog.V(1).InfoS("reconcile for scheduler", "namespace", req.Namespace, "name", req.Name)

	deployable := &appsv1alpha1.Deployable{}
	if err := s.Client.Get(ctx, req.NamespacedName, deployable); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !deployable.ObjectMeta.DeletionTimestamp.IsZero() {
		klog.InfoS("deployable is deleted", "namespace", req.Namespace, "name", req.Name)
		return reconcile.Result{}, nil
	}

	if deployable.Status.PlacementDecided {
		klog.V(1).InfoS("deployable is scheduled, skip", "namespace", deployable.Namespace, "name", deployable.Name)
		return reconcile.Result{}, nil
	}

	return s.scheduleOne(ctx, deployable)
}

func (s *Scheduler) scheduleOne(ctx context.Context, deployable *appsv1alpha1.Deployable) (reconcile.Result, error) {
	deployable.Status.PlacementDecided = true
	deployable.Status.Applied = false

	resNum := len(deployable.Spec.Resources)
	clusterNum := len(deployable.Spec.Placement.ClusterNames)
	if clusterNum == 0 || resNum == 0 {
		return reconcile.Result{}, nil
	}

	deployable.Status.PlacementDecisions = make([]appsv1alpha1.PlacementDecision, clusterNum)
	for i := 0; i < clusterNum; i++ {
		deployable.Status.PlacementDecisions[i] = appsv1alpha1.PlacementDecision{
			Cluster:   deployable.Spec.Placement.ClusterNames[i],
			Resources: deployable.Spec.Resources,
		}
	}

	// bind
	runtimeObject := deployable.DeepCopy()
	_, err := controllerutil.CreateOrPatch(ctx, s.Client, runtimeObject, func() error {
		runtimeObject.Status.Applied = deployable.Status.Applied
		runtimeObject.Status.PlacementDecided = deployable.Status.PlacementDecided
		runtimeObject.Status.PlacementDecisions = deployable.Status.PlacementDecisions
		return nil
	})
	if err != nil {
		klog.ErrorS(err, "unable to create or update for Deployable", "namespace", deployable.Namespace, "name", deployable.Name)
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
