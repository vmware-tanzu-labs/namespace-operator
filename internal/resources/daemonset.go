// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	appsv1 "k8s.io/api/apps/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	DaemonSetKind = "DaemonSet"
)

// DaemonSetIsReady checks to see if a daemonset is ready.
func DaemonSetIsReady(resource common.ComponentResource) (bool, error) {
	var daemonSet appsv1.DaemonSet
	if err := getObject(resource, &daemonSet, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if daemonSet.Name == "" {
		return false, nil
	}

	// ensure the desired number is scheduled and ready
	if daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady {
		if daemonSet.Status.NumberReady > 0 && daemonSet.Status.NumberUnavailable < 1 {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, nil
}
