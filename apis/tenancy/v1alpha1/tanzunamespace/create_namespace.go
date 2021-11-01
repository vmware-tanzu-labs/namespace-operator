// Copyright 2006-2021 VMware, Inc.
// SPDX-License-Identifier: MIT
/*

 */

package resources

import (
	"text/template"

	k8s_api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"

	tenancyv1alpha1 "github.com/vmware-tanzu-labs/namespace-operator/apis/tenancy/v1alpha1"
)

const resourceNamespace = `
---
apiVersion: v1
kind: Namespace
metadata:
{{ $defaultNamespaceName := defaultNamespace .Name .Spec }}
  name: {{ $defaultNamespaceName }}
  labels:
    name: {{ $defaultNamespaceName }}
`

// CreateNamespace creates Namespace resource for the parent TanzuNamespace object
func CreateNamespace(parent *tenancyv1alpha1.TanzuNamespace) (metav1.Object, error) {

	fmap := template.FuncMap{
		"defaultNamespace": defaultNamespace,
	}

	childContent, err := runTemplate("tanzu-namespace", resourceNamespace, parent, fmap)
	if err != nil {
		return nil, err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(childContent), nil, nil)
	if err != nil {
		return nil, err
	}

	resourceObj := obj.(*k8s_api.Namespace)

	return resourceObj, nil
}
