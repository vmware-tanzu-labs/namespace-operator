// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	ServiceKind = "Service"
)

// ServiceIsReady checks to see if a job is ready.
func ServiceIsReady(resource common.ComponentResource) (bool, error) {
	var service corev1.Service
	if err := getObject(resource, &service, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if service.Name == "" {
		return false, nil
	}

	// return if we have an external service type
	if service.Spec.Type == corev1.ServiceTypeExternalName {
		return true, nil
	}

	// ensure a cluster ip address exists for cluster ip types
	if service.Spec.ClusterIP != corev1.ClusterIPNone && len(service.Spec.ClusterIP) == 0 {
		return false, nil
	}

	// ensure a load balancer ip or hostname is present
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) == 0 {
			return false, nil
		}
	}

	return true, nil
}
