/*
Copyright 2017 The Caicloud Authors.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=resourceclass

// ResourceClass is a specification for a ResourceClass resource
type ResourceClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Required. Spec defines resources
	Spec ResourceClassSpec `json:"spec"`
	// Required.
	Status ResourceClassStatus `json:"status"`
}

// ResourceClassSpec is the spec for a ResourceClass resource
type ResourceClassSpec struct {
	// Required. ResourceName is advertised by device plugin.
	ResourceName string `json:"resourceName"`
	// Required. DevicesIDs is advertised by device plugin, we cached here.
	DevicesIDs []string `json:"deviceIDs"`
	// Required. A list of resource selector requirements.
	// The requirements are ANDed.
	MatchExpressions []ResourceSelectorRequirement
	// Required. An array of available nodes for this resource class.
	// If empty then this resource can't be scheduled.
	Nodes []NodeSpec `json:"nodes"`
}

// A resource selector requirement is a selector that contains values, a key, and an operator
// that relates the key and values.
type NodeSpec struct {
	// Name is the name of this node.
	// We will add label to this node with "ResourceClassName=true".
	Name string `json:"name"`
	// Dedicated represents whether this node is dedicated.
	// If true, we will add taint to this node with "ResourceClassName=true:NoSchedule".
	Dedicated bool `json:"dedicated"`
}

// A resource selector requirement is a selector that contains values, a key, and an operator
// that relates the key and values.
type ResourceSelectorRequirement struct {
	// The label key that the selector applies to.
	Key string
	// Represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	Operator ResourceSelectorOperator
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. If the operator is Gt or Lt, the values
	// array must have a single element, which will be interpreted as an integer.
	// This array is replaced during a strategic merge patch.
	// +optional
	Values []string
}

// A resource selector operator is the set of operators that can be used in
// a resource selector requirement.
type ResourceSelectorOperator string

const (
	ResourceSelectorOpIn           ResourceSelectorOperator = "In"
	ResourceSelectorOpNotIn        ResourceSelectorOperator = "NotIn"
	ResourceSelectorOpExists       ResourceSelectorOperator = "Exists"
	ResourceSelectorOpDoesNotExist ResourceSelectorOperator = "DoesNotExist"
	ResourceSelectorOpGt           ResourceSelectorOperator = "Gt"
	ResourceSelectorOpLt           ResourceSelectorOperator = "Lt"
)

// ResourceClassStatus is the status for a ResourceClass resource
type ResourceClassStatus struct {
	// Required. default is true, If false then this resource can't be required.
	Enable bool `json:"enable"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=resourceclass

// ResourceClassList is a list of ResourceClass resources
type ResourceClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ResourceClass `json:"items"`
}
