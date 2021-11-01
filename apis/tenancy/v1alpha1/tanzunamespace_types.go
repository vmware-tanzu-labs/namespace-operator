// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TanzuNamespaceSpec defines the desired state of TanzuNamespace
type TanzuNamespaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Optional
	Name string `json:"name"`

	//
	//
	// BACKWARDS COMPATIBILITY ONLY
	// the below is for backwards compatibility only, please use Name
	//
	//
	// +kubebuilder:validation:Optional
	TanzuNamespaceName string `json:"tanzuNamespaceName"`

	// +kubebuilder:default={}
	// +kubebuilder:validation:Optional
	LimitRange LimitRange `json:"limitRange"`

	//
	//
	// BACKWARDS COMPATIBILITY ONLY
	// the below is for backwards compatibility only, please use LimitRange map instead
	//
	//

	// +kubebuilder:default="125m"
	// +kubebuilder:validation:Optional
	TanzuLimitRangeDefaultCpuLimit string `json:"tanzuLimitRangeDefaultCpuLimit"`

	// +kubebuilder:default="64Mi"
	// +kubebuilder:validation:Optional
	TanzuLimitRangeDefaultMemoryLimit string `json:"tanzuLimitRangeDefaultMemoryLimit"`

	// +kubebuilder:default="125m"
	// +kubebuilder:validation:Optional
	TanzuLimitRangeDefaultCpuRequest string `json:"tanzuLimitRangeDefaultCpuRequest"`

	// +kubebuilder:default="64Mi"
	// +kubebuilder:validation:Optional
	TanzuLimitRangeDefaultMemoryRequest string `json:"tanzuLimitRangeDefaultMemoryRequest"`

	// +kubebuilder:default="1000m"
	// +kubebuilder:validation:Optional
	TanzuLimitRangeMaxCpuLimit string `json:"tanzuLimitRangeMaxCpuLimit"`

	// +kubebuilder:default="2Gi"
	// +kubebuilder:validation:Optional
	TanzuLimitRangeMaxMemoryLimit string `json:"tanzuLimitRangeMaxMemoryLimit"`

	//
	//
	// BACKWARDS COMPATIBILITY ONLY
	// the below is for backwards compatibility only, please use ResourceQuota map instead
	//
	//

	// +kubebuilder:default="2000m"
	// +kubebuilder:validation:Optional
	TanzuResourceQuotaCpuRequests string `json:"tanzuResourceQuotaCpuRequests"`

	// +kubebuilder:default="4Gi"
	// +kubebuilder:validation:Optional
	TanzuResourceQuotaMemoryRequests string `json:"tanzuResourceQuotaMemoryRequests"`

	// +kubebuilder:default="2000m"
	// +kubebuilder:validation:Optional
	TanzuResourceQuotaCpuLimits string `json:"tanzuResourceQuotaCpuLimits"`

	// +kubebuilder:default="4Gi"
	// +kubebuilder:validation:Optional
	TanzuResourceQuotaMemoryLimits string `json:"tanzuResourceQuotaMemoryLimits"`

	// +kubebuilder:default={}
	// +kubebuilder:validation:Optional
	ResourceQuota ResourceQuota `json:"resourceQuota"`

	// +kubebuilder:default={}
	// +kubebuilder:validation:Optional
	NetworkPolicies []NetworkPolicy `json:"networkPolicies"`

	// +kubebuilder:default={}
	// +kubebuilder:validation:Optional
	RBAC []RBAC `json:"rbac"`
}

// LimitRange defines an individual limit range which belongs to a
// TanzuNamespace parent object
type LimitRange struct {
	DefaultCPULimit      string `json:"defaultCPULimit"`
	DefaultMemoryLimit   string `json:"defaultMemoryLimit"`
	DefaultCPURequest    string `json:"defaultCPURequest"`
	DefaultMemoryRequest string `json:"defaultMemoryRequest"`
	MaxCPULimit          string `json:"maxCPULimit"`
	MaxMemoryLimit       string `json:"maxMemoryLimit"`
}

// ResourceQuota defines an individual resource quota which belongs to a
// TanzuNamespace parent object
type ResourceQuota struct {
	RequestsCPU    string `json:"requestsCPU"`
	RequestsMemory string `json:"requestsMemory"`
	LimitsCPU      string `json:"limitsCPU"`
	LimitsMemory   string `json:"limitsMemory"`
}

// NetworkPolicy defines an individual network policy which belongs to
// an array of NetworkPolicies
type NetworkPolicy struct {
	TargetPodLabels        map[string]string `json:"targetPodLabels"`
	IngressNamespaceLabels map[string]string `json:"ingressNamespaceLabels"`
	IngressPodLabels       map[string]string `json:"ingressPodLabels"`
	IngressTCPPorts        []int             `json:"ingressTCPPorts"`
	IngressUDPPorts        []int             `json:"ingressUDPPorts"`
	EgressNamespaceLabels  map[string]string `json:"egressNamespaceLabels"`
	EgressPodLabels        map[string]string `json:"egressPodLabels"`
	EgressTCPPorts         []int             `json:"egressTCPPorts"`
	EgressUDPPorts         []int             `json:"egressUDPPorts"`
}

// NetworkPort defines an individual network port
type NetworkPort struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

// RBAC defines default RBAC settings
type RBAC struct {
	Create      bool   `json:"create"`
	Type        string `json:"type"`
	User        string `json:"user"`
	Role        string `json:"role"`
	RoleBinding string `json:"roleBinding"`
	Permissions string `json:"permissions"`
	Namespace   string `json:"namespace"`
}

// TanzuNamespaceStatus defines the observed state of TanzuNamespace
type TanzuNamespaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Created    bool        `json:"created,omitempty"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// Condition sets the status.conditions field on the object
type Condition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// TanzuNamespace is the Schema for the tanzunamespaces API
type TanzuNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TanzuNamespaceSpec   `json:"spec,omitempty"`
	Status            TanzuNamespaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TanzuNamespaceList contains a list of TanzuNamespace
type TanzuNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TanzuNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TanzuNamespace{}, &TanzuNamespaceList{})
}
