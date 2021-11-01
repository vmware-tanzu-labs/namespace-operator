// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package tanzunamespace

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tenancyv1alpha2 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha2"
)

// CreateFuncs is an array of functions that are called to create the child resources for the controller
// in memory during the reconciliation loop prior to persisting the changes or updates to the Kubernetes
// database.
var CreateFuncs = []func(
	*tenancyv1alpha2.TanzuNamespace) (metav1.Object, error){
	CreateNamespaceParentSpecNamespace,
	CreateLimitRangeTanzuLimitRange,
	CreateResourceQuotaTanzuResourceQuota,
	CreateNetworkPolicyTanzuNetworkPolicy,
}

// InitFuncs is an array of functions that are called prior to starting the controller manager.  This is
// necessary in instances which the controller needs to "own" objects which depend on resources to
// pre-exist in the cluster. A common use case for this is the need to own a custom resource.
// If the controller needs to own a custom resource type, the CRD that defines it must
// first exist. In this case, the InitFunc will create the CRD so that the controller
// can own custom resources of that type.  Without the InitFunc the controller will
// crash loop because when it tries to own a non-existent resource type during manager
// setup, it will fail.
var InitFuncs = []func(
	*tenancyv1alpha2.TanzuNamespace) (metav1.Object, error){}
