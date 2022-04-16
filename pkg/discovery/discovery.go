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

package discovery

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/klog/v2"
	"open-cluster-management.io/api/work/v1"

	"github.com/multi-cluster-platform/mcp/pkg/apis/apps/v1alpha1"
)

// CRDsInstalled checks if the CRDs are installed or not
func CRDsInstalled(discovery *discovery.DiscoveryClient) bool {
	gvs := []schema.GroupVersionKind{
		v1alpha1.SchemeGroupVersion.WithKind("Manifest"),
		v1alpha1.SchemeGroupVersion.WithKind("Deployable"),
		v1.GroupVersion.WithKind("ManifestWork"),
	}

	for _, gv := range gvs {
		if !isCRDInstalled(discovery, gv) {
			return false
		}
	}

	return true
}

func isCRDInstalled(discovery *discovery.DiscoveryClient, gvk schema.GroupVersionKind) bool {
	crdList, err := discovery.ServerResourcesForGroupVersion(gvk.GroupVersion().String())
	if err != nil {
		klog.ErrorS(err, "resource not found", "resource", gvk)
		return false
	}

	for _, crd := range crdList.APIResources {
		if crd.Kind == gvk.Kind {
			klog.InfoS("resource CRD not found", "resource", crd.Kind)
			return true
		}
	}
	return false
}
