// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// CreateResourcesPhase.DefaultRequeue executes checking for a parent components readiness status.
func (phase *CreateResourcesPhase) DefaultRequeue() ctrl.Result {
	return Requeue()
}

// createResourcePhases defines the phases for resource creation and the order in which they run during the reconcile process.
func createResourcePhases() []ResourcePhase {
	return []ResourcePhase{
		// wait for other resources before attempting to create
		&WaitForResourcePhase{},

		// create the resource in the cluster
		&PersistResourcePhase{},
	}
}

// CreateResourcesPhase.Execute executes executes sub-phases which are required to create the resources.
func (phase *CreateResourcesPhase) Execute(
	r common.ComponentReconciler,
) (proceedToNextPhase bool, err error) {
	// execute the resource phases against each resource
	for _, resource := range r.GetResources() {
		resourceCommon := resource.ToCommonResource()
		resourceCondition := &common.ResourceCondition{}

		for _, resourcePhase := range createResourcePhases() {
			r.GetLogger().V(7).Info(fmt.Sprintf("enter resource phase: %T", resourcePhase))
			_, proceed, err := resourcePhase.Execute(resource, *resourceCondition)

			// set a message, return the error and result on error or when unable to proceed
			if err != nil || !proceed {
				return handleResourcePhaseExit(r, *resourceCommon, *resourceCondition, resourcePhase, proceed, err)
			}

			// set attributes on the resource condition before updating the status
			resourceCondition.LastResourcePhase = getResourcePhaseName(resourcePhase)

			r.GetLogger().V(5).Info(fmt.Sprintf("completed resource phase: %T", resourcePhase))
		}
	}

	return true, nil
}
