// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	JobKind = "Job"
)

// JobIsReady checks to see if a job is ready.
func JobIsReady(resource common.ComponentResource) (bool, error) {
	var job batchv1.Job
	if err := getObject(resource, &job, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if job.Name == "" {
		return false, nil
	}

	// return immediately if the job is active or has no completion time
	if job.Status.Active == 1 || job.Status.CompletionTime == nil {
		return false, nil
	}

	// ensure the completion is actually successful
	if job.Status.Succeeded != 1 {
		return false, fmt.Errorf("job " + job.GetName() + " was not successful")
	}

	return true, nil
}
