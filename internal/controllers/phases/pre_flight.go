// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package phases

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// PreFlightPhase.DefaultRequeue executes checking for a parent components readiness status.
func (phase *PreFlightPhase) DefaultRequeue() ctrl.Result {
	return Requeue()
}

// PreFlightPhase.Execute executes pre-flight and fail-fast conditions prior to attempting resource creation.
func (phase *PreFlightPhase) Execute(
	r common.ComponentReconciler,
) (proceedToNextPhase bool, err error) {
	return true, nil
}
