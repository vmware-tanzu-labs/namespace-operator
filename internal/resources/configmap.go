// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	v1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	ConfigMapKind = "ConfigMap"
)

// ConfigMapIsReady performs the logic to determine if a secret is ready.
func ConfigMapIsReady(resource common.ComponentResource, expectedKeys ...string) (bool, error) {
	var configMap v1.ConfigMap
	if err := getObject(resource, &configMap, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if configMap.Name == "" {
		return false, nil
	}

	// check the status for a ready ca keypair secret
	for _, key := range expectedKeys {
		if string(configMap.Data[key]) == "" {
			return false, nil
		}
	}

	return true, nil
}
