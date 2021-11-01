// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package resources

import (
	"k8s.io/apimachinery/pkg/api/errors"

	extensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/vmware-tanzu-labs/namespace-operator/apis/common"
)

const (
	CustomResourceDefinitionKind = "CustomResourceDefinition"
)

// CustomResourceDefinitionIsReady performs the logic to determine if a custom resource definition is ready.
func CustomResourceDefinitionIsReady(resource common.ComponentResource) (bool, error) {
	var crd extensionsv1.CustomResourceDefinition
	if err := getObject(resource, &crd, false); err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
	}

	return true, nil
}
