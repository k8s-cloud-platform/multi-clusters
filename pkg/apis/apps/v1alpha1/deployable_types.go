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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=deployables,scope=Namespaced,categories=mcp-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Deployable is the deploy unit for Manifests, use mcp-system as namespace for cluster scope resource
type Deployable struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DeployableSpec `json:"spec"`

	// +optional
	Status DeployableStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeployableList contains a list of Deployable
type DeployableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deployable `json:"items"`
}

type DeployableSpec struct {
	// +optional
	Placement Placement `json:"placement,omitempty"`

	// +optional
	Resources []corev1.ObjectReference `json:"resources,omitempty"`
}

type Placement struct {
	// +optional
	ClusterNames []string `json:"clusterNames,omitempty"`
}

type DeployableStatus struct {
	// scheduler handled
	// +optional
	PlacementDecided bool `json:"placementDecided"`

	// ManifestWork generated
	// +optional
	Applied bool `json:"applied"`

	// +optional
	PlacementDecisions []PlacementDecision `json:"placementDecisions,omitempty"`
}

type PlacementDecision struct {
	// +optional
	Cluster string `json:"cluster,omitempty"`

	// TODO, maybe add override to each resource
	// +optional
	Resources []corev1.ObjectReference `json:"resources,omitempty"`
}
