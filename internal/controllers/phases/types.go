// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package phases

import (
	"fmt"
	"reflect"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// Phase defines a phase of the reconciliation process.
type Phase interface {
	Execute(common.ComponentReconciler) (bool, error)
	DefaultRequeue() ctrl.Result
}

// ResourcePhase defines the specific phase of reconcilication associated with creating resources.
type ResourcePhase interface {
	Execute(common.ComponentResource, common.ResourceCondition) (ctrl.Result, bool, error)
}

// Below are the phase types which satisfy the Phase interface.
type DependencyPhase struct{}
type PreFlightPhase struct{}
type CreateResourcesPhase struct{}
type CheckReadyPhase struct{}
type CompletePhase struct{}

// Below are the phase types which satisfy the ResourcePhase interface.
type PersistResourcePhase struct{}
type WaitForResourcePhase struct{}

// GetSuccessCondition defines the success condition for the phase.
func GetSuccessCondition(phase Phase) common.PhaseCondition {
	return common.PhaseCondition{
		Phase:   getPhaseName(phase),
		State:   common.PhaseStateComplete,
		Message: "Successfully Completed Phase",
	}
}

// GetPendingCondition defines the pending condition for the phase.
func GetPendingCondition(phase Phase) common.PhaseCondition {
	return common.PhaseCondition{
		Phase:   getPhaseName(phase),
		State:   common.PhaseStatePending,
		Message: "Pending Execution of Phase",
	}
}

// GetFailCondition defines the fail condition for the phase.
func GetFailCondition(phase Phase, err error) common.PhaseCondition {
	return common.PhaseCondition{
		Phase:   getPhaseName(phase),
		State:   common.PhaseStateFailed,
		Message: "Failed Phase with Error; " + err.Error(),
	}
}

func getPhaseName(phase Phase) string {
	objectElements := strings.Split(fmt.Sprintf("%s", reflect.TypeOf(phase)), ".")

	return objectElements[len(objectElements)-1]
}

func getResourcePhaseName(resourcePhase ResourcePhase) string {
	objectElements := strings.Split(fmt.Sprintf("%s", reflect.TypeOf(resourcePhase)), ".")

	return objectElements[len(objectElements)-1]
}
