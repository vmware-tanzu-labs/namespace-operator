// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package phases

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// CompletePhase.DefaultRequeue executes checking for a parent components readiness status.
func (phase *CompletePhase) DefaultRequeue() ctrl.Result {
	return Requeue()
}

// CompletePhase.Execute executes the completion of a reconciliation loop.
func (phase *CompletePhase) Execute(
	r common.ComponentReconciler,
) (proceedToNextPhase bool, err error) {
	r.GetComponent().SetReadyStatus(true)
	r.GetLogger().V(0).Info("successfully reconciled")

	return true, nil
}
