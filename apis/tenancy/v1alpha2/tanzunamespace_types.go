// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package v1alpha2

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TanzuNamespaceSpec defines the desired state of TanzuNamespace.
type TanzuNamespaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Namespace name which is created and then enforced by related policy objects
	// such as LimitRange, ResourceQuota, and NetworkPolicy.
	Namespace string `json:"namespace"`

	Resources TanzuNamespaceSpecResources `json:"resources"`
}

type TanzuNamespaceSpecResources struct {
	Limits TanzuNamespaceSpecResourcesLimits `json:"limits"`

	Requests TanzuNamespaceSpecResourcesRequests `json:"requests"`

	Max TanzuNamespaceSpecResourcesMax `json:"max"`

	Quota TanzuNamespaceSpecResourcesQuota `json:"quota"`
}

type TanzuNamespaceSpecResourcesLimits struct {
	// +kubebuilder:default="250m"
	// +kubebuilder:validation:Optional
	// Default CPU limits to be applied to applications which get deployed into this namespace,
	// but are missing a resources declaration.
	Cpu string `json:"cpu"`

	// +kubebuilder:default="64Mi"
	// +kubebuilder:validation:Optional
	// Default Memory limits to be applied to applications which get deployed into this namespace,
	// but are missing a resources declaration.
	Memory string `json:"memory"`
}

type TanzuNamespaceSpecResourcesRequests struct {
	// +kubebuilder:default="250m"
	// +kubebuilder:validation:Optional
	// Default CPU requests to be applied to applications which get deployed into this namespace,
	// but are missing a resources declaration.
	Cpu string `json:"cpu"`

	// +kubebuilder:default="64Mi"
	// +kubebuilder:validation:Optional
	// Default Memory requests to be applied to applications which get deployed into this namespace,
	// but are missing a resources declaration.
	Memory string `json:"memory"`
}

type TanzuNamespaceSpecResourcesMax struct {
	// +kubebuilder:default="500m"
	// +kubebuilder:validation:Optional
	// Default maximum CPU limits for an individual application which get deployed into this namespace.
	Cpu string `json:"cpu"`

	// +kubebuilder:default="256Mi"
	// +kubebuilder:validation:Optional
	// Default maximum Memory limits for an individual application which get deployed into this namespace.
	Memory string `json:"memory"`
}

type TanzuNamespaceSpecResourcesQuota struct {
	Requests TanzuNamespaceSpecResourcesQuotaRequests `json:"requests"`

	Limits TanzuNamespaceSpecResourcesQuotaLimits `json:"limits"`
}

type TanzuNamespaceSpecResourcesQuotaRequests struct {
	// +kubebuilder:default="2000m"
	// +kubebuilder:validation:Optional
	// Default CPU requests quota to be enforced on the sum of all applications which get deployed into this namespace.
	Cpu string `json:"cpu"`

	// +kubebuilder:default="4Gi"
	// +kubebuilder:validation:Optional
	// Default Memory requests quota to be enforced on the sum of all applications which get deployed into this namespace.
	Memory string `json:"memory"`
}

type TanzuNamespaceSpecResourcesQuotaLimits struct {
	// +kubebuilder:default="2000m"
	// +kubebuilder:validation:Optional
	// Default CPU limits quota to be enforced on the sum of all applications which get deployed into this namespace.
	Cpu string `json:"cpu"`

	// +kubebuilder:default="4Gi"
	// +kubebuilder:validation:Optional
	// Default Memory limits quota to be enforced on the sum of all applications which get deployed into this namespace.
	Memory string `json:"memory"`
}

// TanzuNamespaceStatus defines the observed state of TanzuNamespace.
type TanzuNamespaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Created               bool                    `json:"created,omitempty"`
	DependenciesSatisfied bool                    `json:"dependenciesSatisfied,omitempty"`
	Conditions            []common.PhaseCondition `json:"conditions,omitempty"`
	Resources             []common.Resource       `json:"resources,omitempty"`
}

// +kubebuilder:storageversion
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// TanzuNamespace is the Schema for the tanzunamespaces API.
type TanzuNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TanzuNamespaceSpec   `json:"spec,omitempty"`
	Status            TanzuNamespaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TanzuNamespaceList contains a list of TanzuNamespace.
type TanzuNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TanzuNamespace `json:"items"`
}

// interface methods

// GetReadyStatus returns the ready status for a component.
func (component *TanzuNamespace) GetReadyStatus() bool {
	return component.Status.Created
}

// SetReadyStatus sets the ready status for a component.
func (component *TanzuNamespace) SetReadyStatus(status bool) {
	component.Status.Created = status
}

// GetDependencyStatus returns the dependency status for a component.
func (component *TanzuNamespace) GetDependencyStatus() bool {
	return component.Status.DependenciesSatisfied
}

// SetDependencyStatus sets the dependency status for a component.
func (component *TanzuNamespace) SetDependencyStatus(dependencyStatus bool) {
	component.Status.DependenciesSatisfied = dependencyStatus
}

// GetPhaseConditions returns the phase conditions for a component.
func (component TanzuNamespace) GetPhaseConditions() []common.PhaseCondition {
	return component.Status.Conditions
}

// SetPhaseCondition sets the phase conditions for a component.
func (component *TanzuNamespace) SetPhaseCondition(condition common.PhaseCondition) {
	if found := condition.GetPhaseConditionIndex(component); found >= 0 {
		if condition.LastModified == "" {
			condition.LastModified = time.Now().UTC().String()
		}
		component.Status.Conditions[found] = condition
	} else {
		component.Status.Conditions = append(component.Status.Conditions, condition)
	}
}

// GetResources returns the resources for a component.
func (component TanzuNamespace) GetResources() []common.Resource {
	return component.Status.Resources
}

// SetResources sets the phase conditions for a component.
func (component *TanzuNamespace) SetResource(resource common.Resource) {

	if found := resource.GetResourceIndex(component); found >= 0 {
		if resource.ResourceCondition.LastModified == "" {
			resource.ResourceCondition.LastModified = time.Now().UTC().String()
		}
		component.Status.Resources[found] = resource
	} else {
		component.Status.Resources = append(component.Status.Resources, resource)
	}
}

// GetDependencies returns the dependencies for a component.
func (*TanzuNamespace) GetDependencies() []common.Component {
	return []common.Component{}
}

// GetComponentGVK returns a GVK object for the component.
func (*TanzuNamespace) GetComponentGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		Kind:    "TanzuNamespace",
	}
}

func init() {
	SchemeBuilder.Register(&TanzuNamespace{}, &TanzuNamespaceList{})
}
