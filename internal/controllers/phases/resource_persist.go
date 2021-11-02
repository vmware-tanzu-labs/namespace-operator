// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package phases

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

// PersistResourcePhase.Execute executes persisting resources to the Kubernetes database.
func (phase *PersistResourcePhase) Execute(
	resource common.ComponentResource,
	resourceCondition common.ResourceCondition,
) (ctrl.Result, bool, error) {
	// persist the resource
	if err := persistResource(
		resource,
		resourceCondition,
		phase,
	); err != nil {
		return ctrl.Result{}, false, err
	}

	return ctrl.Result{}, true, nil
}

// persistResource persists a single resource to the Kubernetes database.
func persistResource(
	resource common.ComponentResource,
	condition common.ResourceCondition,
	phase *PersistResourcePhase,
) error {
	// persist resource
	r := resource.GetReconciler()
	if err := r.CreateOrUpdate(resource.GetObject()); err != nil {
		if IsOptimisticLockError(err) {
			return nil
		} else {
			r.GetLogger().V(0).Info(err.Error())

			return err
		}
	}

	// set attributes related to the persistence of this child resource
	condition.LastResourcePhase = getResourcePhaseName(phase)
	condition.LastModified = time.Now().UTC().String()
	condition.Message = "resource created successfully"
	condition.Created = true

	// update the condition to notify that we have created a child resource
	return updateResourceConditions(r, *resource.ToCommonResource(), &condition)
}
