// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	appsv1 "k8s.io/api/apps/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	DeploymentKind = "Deployment"
)

// DeploymentIsReady performs the logic to determine if a deployment is ready.
func DeploymentIsReady(resource common.ComponentResource) (bool, error) {
	var deployment appsv1.Deployment
	if err := getObject(resource, &deployment, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if deployment.Name == "" {
		return false, nil
	}

	// rely on observed generation to give us a proper status
	if deployment.Generation != deployment.Status.ObservedGeneration {
		return false, nil
	}

	// check the status for a ready deployment
	if deployment.Status.ReadyReplicas != deployment.Status.Replicas {
		return false, nil
	}

	return true, nil
}
