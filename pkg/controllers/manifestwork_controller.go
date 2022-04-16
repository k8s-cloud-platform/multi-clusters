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

package controllers

import (
	"context"
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	workv1 "open-cluster-management.io/api/work/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "github.com/multi-cluster-platform/mcp/pkg/apis/apps/v1alpha1"
	"github.com/multi-cluster-platform/mcp/pkg/constants"
)

type ManifestWorkController struct {
	client.Client
	client.Reader
}

var _ reconcile.Reconciler = &ManifestWorkController{}

// SetupWithManager sets up the controller with the Manager.
func (c *ManifestWorkController) SetupWithManager(mgr ctrl.Manager, options controller.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Deployable{}).
		WithOptions(options).
		Complete(c)
}

func (c *ManifestWorkController) Reconcile(ctx context.Context, req reconcile.Request) (_ reconcile.Result, reterr error) {
	klog.V(1).InfoS("reconcile for Deployable", "namespace", req.Namespace, "name", req.Name)

	deployable := &appsv1alpha1.Deployable{}
	if err := c.Client.Get(ctx, req.NamespacedName, deployable); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	defer func() {
		runtimeObject := deployable.DeepCopy()
		_, err := controllerutil.CreateOrPatch(ctx, c.Client, runtimeObject, func() error {
			runtimeObject.ObjectMeta.Finalizers = deployable.ObjectMeta.Finalizers
			runtimeObject.Status.Applied = deployable.Status.Applied
			return nil
		})
		if err != nil {
			klog.ErrorS(err, "unable to create or patch Deployable", "namespace", deployable.Namespace, "name", deployable.Name)
			reterr = kerrors.NewAggregate([]error{reterr, err})
		}
	}()

	// Add finalizer first if not exist to avoid the race condition between init and delete
	if !controllerutil.ContainsFinalizer(deployable, constants.DeployableFinalizer) {
		controllerutil.AddFinalizer(deployable, constants.DeployableFinalizer)
		return ctrl.Result{}, nil
	}

	if !deployable.ObjectMeta.DeletionTimestamp.IsZero() {
		return c.reconcileDelete(ctx, deployable)
	}

	return c.reconcileNormal(ctx, deployable)
}

func (c *ManifestWorkController) reconcileDelete(ctx context.Context, deployable *appsv1alpha1.Deployable) (reconcile.Result, error) {
	klog.V(1).InfoS("reconcile for Deployable delete", "namespace", deployable.Namespace, "name", deployable.Name)

	if deployable.Status.PlacementDecided {
		// delete decided resources
		for _, decision := range deployable.Status.PlacementDecisions {
			manifestWork := &workv1.ManifestWork{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: decision.Cluster,
					Name:      fmt.Sprintf("%s-%s", deployable.Namespace, deployable.Name),
				},
			}
			if err := c.Client.Delete(ctx, manifestWork); err != nil {
				if apierrors.IsNotFound(err) {
					continue
				}
				return reconcile.Result{}, err
			}
		}
	}

	deployable.Status.Applied = false
	controllerutil.RemoveFinalizer(deployable, constants.DeployableFinalizer)
	return reconcile.Result{}, nil
}

func (c *ManifestWorkController) reconcileNormal(ctx context.Context, deployable *appsv1alpha1.Deployable) (reconcile.Result, error) {
	klog.V(1).InfoS("reconcile for Deployable normal", "namespace", deployable.Namespace, "name", deployable.Name)

	if !deployable.Status.PlacementDecided {
		klog.V(1).InfoS("deployable is not scheduled, skip", "namespace", deployable.Namespace, "name", deployable.Name)
		return reconcile.Result{}, nil
	}

	if deployable.Status.Applied {
		klog.V(1).InfoS("deployable is already applied, skip", "namespace", deployable.Namespace, "name", deployable.Name)
		return reconcile.Result{}, nil
	}

	for _, decision := range deployable.Status.PlacementDecisions {
		// handle for each ManifestWork
		manifestWork := &workv1.ManifestWork{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: decision.Cluster,
				Name:      fmt.Sprintf("%s-%s", deployable.Namespace, deployable.Name),
			},
			Spec: workv1.ManifestWorkSpec{
				Workload: workv1.ManifestsTemplate{
					Manifests: make([]workv1.Manifest, len(decision.Resources)),
				},
			},
		}

		for idx, resource := range decision.Resources {
			manifest := &appsv1alpha1.Manifest{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: resource.Namespace,
					Name:      strings.ToLower(fmt.Sprintf("%s-%s-%s", convertAPIVersion(resource.APIVersion), resource.Kind, resource.Name)),
				},
			}
			if err := c.Client.Get(ctx, client.ObjectKeyFromObject(manifest), manifest); err != nil {
				klog.ErrorS(err, "unable to get Manifest", "namespace", manifestWork.Namespace, "name", manifestWork.Name)
				return reconcile.Result{}, err
			}

			manifestWork.Spec.Workload.Manifests[idx] = workv1.Manifest{
				RawExtension: manifest.Template,
			}
		}

		runtimeObject := manifestWork.DeepCopy()
		result, err := controllerutil.CreateOrUpdate(ctx, c.Client, runtimeObject, func() error {
			runtimeObject.Spec = manifestWork.Spec
			return nil
		})
		if err != nil {
			klog.ErrorS(err, "unable to create or update for ManifestWork", "namespace", manifestWork.Namespace, "name", manifestWork.Name)
			return reconcile.Result{}, err
		}

		if result == controllerutil.OperationResultCreated {
			klog.V(1).InfoS("success to create ManifestWork", "namespace", manifestWork.Namespace, "name", manifestWork.Name)
		} else if result == controllerutil.OperationResultUpdated {
			klog.V(1).InfoS("success to update ManifestWork", "namespace", manifestWork.Namespace, "name", manifestWork.Name)
		}
	}

	deployable.Status.Applied = true

	return reconcile.Result{}, nil
}

func convertAPIVersion(apiVersion string) string {
	return strings.Join(strings.Split(apiVersion, "/"), "-")
}
