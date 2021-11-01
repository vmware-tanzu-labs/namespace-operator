// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	v1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	SecretKind = "Secret"
)

// SecretIsReady performs the logic to determine if a secret is ready.
func SecretIsReady(resource common.ComponentResource, expectedKeys ...string) (bool, error) {
	var secret v1.Secret
	if err := getObject(resource, &secret, true); err != nil {
		return false, err
	}

	// if we have a name that is empty, we know we did not find the object
	if secret.Name == "" {
		return false, nil
	}

	// check the status for a ready secret if we expect certain fields to exist
	for _, key := range expectedKeys {
		if string(secret.Data[key]) == "" {
			return false, nil
		}
	}

	return true, nil
}
