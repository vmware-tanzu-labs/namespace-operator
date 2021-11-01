// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package phases

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
	"github.com/vmware-tanzu-labs/namespace-operator/internal/resources"
)

// CheckReadyPhase.DefaultRequeue executes checking for a parent components readiness status.
func (phase *CheckReadyPhase) DefaultRequeue() ctrl.Result {
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: 5 * time.Second,
	}
}

// CheckReadyPhase.Execute executes checking for a parent components readiness status.
func (phase *CheckReadyPhase) Execute(
	r common.ComponentReconciler,
) (proceedToNextPhase bool, err error) {
	// check to see if known types are ready
	knownReady, err := resources.AreReady(r.GetResources()...)
	if err != nil {
		return false, err
	}

	// check to see if the custom methods return ready
	customReady, err := r.CheckReady()
	if err != nil {
		return false, err
	}

	return (knownReady && customReady), nil
}
